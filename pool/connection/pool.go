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

func newDefaultClientPool() *ClientPool {
	return &ClientPool{
		Logger:            log.DummyLogger,
		KeepAliveInterval: DefaultKeepAliveInterval,
		CreateNewInterval: DefaultCreateNewInterval,
		DialRetryCount:    DefaultDialRetryCount,
		DialRetryInterval: DefaultRetryInterval,
	}
}

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

	pool := newDefaultClientPool()

	for _, opt := range opts {
		opt.apply(pool)
	}

	if pool.Logger == nil {
		return nil, ErrWrapper(IllegalParams)
	}

	if pool.Dial == nil ||
		pool.Close == nil ||
		pool.KeepAlive == nil {
		pool.Logger.Errorf("illegal params => pool.Dial | pool.Close | pool.KeepAlive")
		return nil, ErrWrapper(IllegalParams)
	}

	if pool.InitialPoolSize < 0 {
		pool.Logger.Errorf("illegal params => pool.InitialPoolSize < 0")
		return nil, ErrWrapper(IllegalParams)
	}

	if pool.MaxPoolSize < 1 {
		pool.Logger.Errorf("illegal params => pool.MaxPoolSize < 1")
		return nil, ErrWrapper(IllegalParams)
	}

	if pool.InitialPoolSize > pool.MaxPoolSize {
		pool.Logger.Errorf("illegal params => pool.InitialPoolSize > pool.MaxPoolSize")
		return nil, ErrWrapper(IllegalParams)
	}

	pool.retryPool = make(chan int, pool.MaxPoolSize)
	pool.alivePool = make(chan interface{}, pool.MaxPoolSize)
	pool.swapPool = make(chan interface{}, pool.MaxPoolSize)

	for i := 0; i < pool.InitialPoolSize; i++ {
		if c, err := pool.Dial(ctx, pool.Address, pool.Port); err == nil {
			pool.alivePool <- c
		} else {
			pool.Logger.Errorf("initialize pool => Dial %v:%v failed,err:%v", pool.Address, pool.Port, err.Error())
			pool.retryPool <- 0
		}
	}

	go pool.start()

	return pool, nil
}

func (p *ClientPool) start() {
	go p.retryLoop()
	go p.keepAliveLoop()
}

func (p *ClientPool) Get() (connection interface{}, err error) {

	select {
	case <-time.After(p.CreateNewInterval):
		p.sync.Lock()
		defer p.sync.Unlock()

		p.Logger.Warnf("addr:%v => Get new connection from new create.", p.Address)
		if int(p.workConnCount)+len(p.retryPool)+len(p.alivePool)+len(p.swapPool) < p.MaxPoolSize {

			retry := 0
			for retry < p.DialRetryCount {
				if connection, err = p.Dial(context.TODO(), p.Address, p.Port); err != nil {
					p.Logger.Errorf("addr:%v => get conn => Dial %v:%v failed,err:%v", p.Address, p.Address, p.Port, err.Error())
					retry++
					continue
				} else {
					atomic.AddInt32(&p.workConnCount, 1)
					p.Logger.Errorf("addr:%v => get conn => Dial %v:%v successfully", p.Address, p.Address, p.Port)
					return
				}
			}

			if retry >= p.DialRetryCount {
				p.retryPool <- 0
				return nil, err
			}
		} else {
			p.Logger.Errorf("addr:%v => Pool Was Exhausted, detail: working: %v, alive: %v, retry: %v.", p.Address, p.workConnCount, len(p.alivePool), len(p.retryPool))
			return nil, ErrWrapper(PoolExhausted)
		}
	case connection = <-p.alivePool:
		p.Logger.Infof("addr:%v => Get new connection from alive pool.", p.Address)
		atomic.AddInt32(&p.workConnCount, 1)
		return
	case connection = <-p.swapPool:
		p.Logger.Infof("addr:%v => Get new connection from swap pool.", p.Address)
		atomic.AddInt32(&p.workConnCount, 1)
		return
	}

	return nil, ErrWrapper(GetConnTimeout)
}

func (p *ClientPool) Put(connection interface{}) (err error) {

	p.sync.Lock()

	if connection != nil {
		if p.isStopped {
			err := p.Close(context.TODO(), connection)
			if err != nil {
				p.Logger.Errorf("addr:%v => Put conn => Close conn failed,err:%v", p.Address, err.Error())
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
			p.Logger.Infof("addr:%v => Release connection error:%v", p.Address, err)
		}
		atomic.SwapInt32(&p.workConnCount, p.workConnCount-1)
	}

	p.sync.Unlock()
}

func (p *ClientPool) retryLoop() {
	p.Logger.Infof("addr:%v => retry loop start.", p.Address)

	for {
		select {
		case <-time.After(p.DialRetryInterval):
			max := len(p.retryPool)
			for i := 0; i < max; i++ {
				if connection, err := p.Dial(context.TODO(), p.Address, p.Port); err == nil {
					<-p.retryPool
					p.alivePool <- connection
					p.Logger.Infof("addr:%v => Retry Pool Success.", p.Address)
				} else {
					p.Logger.Errorw("addr:%v => Retry Pool Failed.", p.Address)
				}
			}

			if p.isStopped {
				break
			}
		}
	}
}

func (p *ClientPool) keepAliveLoop() {

	p.Logger.Infof("addr:%v => keepAlive loop start.", p.Address)

	for {
		select {
		case <-time.After(p.KeepAliveInterval):

			if len(p.alivePool) > 0 {
				// send keep alive message to each connection
				for connection := range p.alivePool {
					if err := p.KeepAlive(context.TODO(), connection); err == nil {
						p.swapPool <- connection
					} else {
						p.Logger.Errorw("addr:%v => Keepalive Pool Failed on %v\n", p.Address, fmt.Sprintf("%v:%v", p.Address, p.Port))
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

	p.Logger.Infof("addr:%v => keepAlive loop end.", p.Address)
}
