package connection

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"sync"
	"sync/atomic"
	"time"
)

var (
	DefaultDialRetryCount    = 3
	DefaultRetryInterval     = time.Second * 10
	DefaultKeepAliveInterval = time.Second * 3
	DefaultCreateNewInterval = time.Second * 1
)


type ClientPool struct {
	Address           string
	Port              string
	Dial              func(ctx context.Context, address, port string) (connection interface{}, err error)
	Close             func(ctx context.Context, connection interface{}) (err error)
	KeepAlive         func(ctx context.Context, connection interface{}) (err error)
	InitialPoolSize   int
	MaxPoolSize       int
	DialRetryCount    int
	KeepAliveInterval time.Duration
	DialRetryInterval time.Duration
	CreateNewInterval time.Duration
	Logger            log.Logger

	workConnCount int32
	alivePool     chan interface{}
	swapPool      chan interface{}
	retryPool     chan int
	sync          sync.Mutex
	isStopped     bool
}

func NewClientPool(ctx context.Context, opts ...Option) (*ClientPool, error) {

	pool := &ClientPool{
		KeepAliveInterval: DefaultKeepAliveInterval,
		CreateNewInterval: DefaultCreateNewInterval,
		DialRetryCount:    DefaultDialRetryCount,
		DialRetryInterval: DefaultRetryInterval,
	}

	for _, opt := range opts {
		opt.apply(pool)
	}

	if pool.Logger == nil {
		return nil, GenErr(IllegalParams)
	}

	if pool.Dial == nil ||
		pool.Close == nil ||
		pool.KeepAlive == nil {
		pool.Logger.Errorf("illegal params => pool.Dial | pool.Close | pool.KeepAlive")
		return nil, GenErr(IllegalParams)
	}

	if pool.InitialPoolSize < 0 {
		pool.Logger.Errorf("illegal params => pool.InitialPoolSize < 0")
		return nil, GenErr(IllegalParams)
	}

	if pool.MaxPoolSize < 1 {
		pool.Logger.Errorf("illegal params => pool.MaxPoolSize < 1")
		return nil, GenErr(IllegalParams)
	}

	if pool.InitialPoolSize > pool.MaxPoolSize {
		pool.Logger.Errorf("illegal params => pool.InitialPoolSize > pool.MaxPoolSize")
		return nil, GenErr(IllegalParams)
	}

	pool.retryPool = make(chan int, pool.MaxPoolSize)
	pool.alivePool = make(chan interface{}, pool.MaxPoolSize)
	pool.swapPool = make(chan interface{}, pool.MaxPoolSize)

	for i := 0; i < pool.InitialPoolSize; i++ {
		if c, err := pool.Dial(ctx, pool.Address, pool.Port); err == nil {
			pool.alivePool <- c
		} else {
			pool.retryPool <- 0
		}
	}

	return pool, nil
}

func (p *ClientPool) Start() {
	go p.retryLoop()
	go p.keepAliveLoop()
}

func (p *ClientPool) Get() (connection interface{}, err error) {

	select {
	case <-time.After(p.CreateNewInterval):
		p.sync.Lock()
		defer p.sync.Unlock()

		p.Logger.Infof("Get new connection from new create.")
		if int(p.workConnCount)+len(p.retryPool)+len(p.alivePool)+len(p.swapPool) < p.MaxPoolSize {

			retry := 0
			for retry < p.DialRetryCount {
				if connection, err = p.Dial(context.TODO(), p.Address, p.Port); err != nil {
					retry++
					continue
				} else {
					atomic.AddInt32(&p.workConnCount, 1)
					return
				}
			}

			if retry >= p.DialRetryCount {
				p.retryPool <- 0
				return nil, err
			}
		} else {
			p.Logger.Errorf("Pool Was Exhausted, detail: working: %v, alive: %v, retry: %v.", p.workConnCount, len(p.alivePool), len(p.retryPool))
			return nil, GenErr(PoolExhausted)
		}
	case connection = <-p.alivePool:
		p.Logger.Infof("Get new connection from alive pool.")
		atomic.AddInt32(&p.workConnCount, 1)
		return
	case connection = <-p.swapPool:
		p.Logger.Infof("Get new connection from swap pool.")
		atomic.AddInt32(&p.workConnCount, 1)
		return
	}

	return nil, GenErr(GetConnTimeout)
}

func (p *ClientPool) Put(connection interface{}) (err error) {

	p.sync.Lock()

	if connection != nil {
		if p.isStopped {
			err := p.Close(context.TODO(), connection)
			if err != nil {
				return err
			}
		} else {
			if len(p.alivePool) < p.MaxPoolSize {
				p.alivePool <- connection
			}
		}
	}

	atomic.SwapInt32(&p.workConnCount, p.workConnCount-1)
	p.sync.Unlock()

	return
}

func (p *ClientPool) Release() {
	p.sync.Lock()
	p.isStopped = true

	for connection := range p.alivePool {
		if err := p.Close(context.TODO(), connection); err != nil {
			p.Logger.Infof("Release connection error: ", err)
		}
		atomic.SwapInt32(&p.workConnCount, p.workConnCount-1)
	}

	p.sync.Unlock()
}

func (p *ClientPool) retryLoop() {
	p.Logger.Infof("retry loop start.")

exit:
	for {
		select {
		case <-time.After(p.DialRetryInterval):
			max := len(p.retryPool)
			for i := 0; i < max; i++ {
				if connection, err := p.Dial(context.TODO(), p.Address, p.Port); err == nil {
					<-p.retryPool
					p.alivePool <- connection
					p.Logger.Infof("Retry Pool Success.")
				} else {
					p.Logger.Infof("Retry Pool Failed.")
				}
			}

			if p.isStopped {
				break exit
			}
		}
	}

	p.Logger.Infof("retry loop end.")
}

func (p *ClientPool) keepAliveLoop() {

	p.Logger.Infof("keepAlive loop start.")

	for {
		select {
		case <-time.After(p.KeepAliveInterval):

			if len(p.alivePool) > 0 {
				// send keep alive message to each connection
				for connection := range p.alivePool {
					if err := p.KeepAlive(context.TODO(), connection); err == nil {
						p.swapPool <- connection
					} else {
						p.Logger.Infof("Keepalive Pool Failed on %v\n", fmt.Sprintf("%v:%v", p.Address, p.Port))
						p.retryPool <- 0
					}

					if len(p.alivePool) == 0 {
						break
					}
				}
			}

			if len(p.swapPool) > 0 {
				// restore alive connection pool.
				for connection := range p.swapPool {
					p.alivePool <- connection

					if len(p.swapPool) == 0 {
						break
					}
				}
			}
		}

		if p.isStopped {
			for connection := range p.alivePool {
				err := p.Close(context.TODO(), connection)
				if err != nil {
					return
				}
			}
			break
		}
	}

	p.Logger.Infof("keepAlive loop end.")
}
