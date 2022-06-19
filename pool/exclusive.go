package pool

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/easylog"
	"sync/atomic"
	"time"
)

func newExclusivePool(cfg *Cfg) *ExclusivePool {
	return &ExclusivePool{
		Address:           cfg.Address,
		Port:              cfg.Port,
		Dial:              cfg.Dial,
		Close:             cfg.Close,
		KeepAlive:         cfg.KeepAlive,
		InitialPoolSize:   cfg.InitialPoolSize,
		BestPoolSize:      cfg.BestPoolSize,
		MaxPoolSize:       cfg.MaxPoolSize,
		DialRetryCount:    cfg.DialRetryCount,
		KeepAliveInterval: cfg.KeepAliveInterval,
		CleanInterval:     cfg.CleanInterval,
		DialRetryInterval: cfg.DialRetryInterval,
		Logger:            cfg.Logger,
	}
}

type ExclusivePool struct {
	Address           string
	Port              string
	Dial              func(ctx context.Context, address, port string) (connection interface{}, err error)
	Close             func(ctx context.Context, connection interface{}) (err error)
	KeepAlive         func(ctx context.Context, connection interface{}) (err error)
	InitialPoolSize   int
	BestPoolSize      int
	MaxPoolSize       int
	DialRetryCount    int
	KeepAliveInterval time.Duration
	CleanInterval     time.Duration
	DialRetryInterval time.Duration
	Logger            easylog.API

	workConnCount  int32
	newlyConnCount int32
	alivePool      chan interface{}
	swapPool       chan interface{}
	retryPool      chan int
	isStopped      bool
}

func NewExclusivePool(ctx context.Context, cfg *Cfg) (*ExclusivePool, error) {

	pool := newExclusivePool(cfg)

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

	if pool.BestPoolSize == 0 {
		pool.Logger.Errorf("initialize pool => BestPoolSize is not set,use MaxPoolSize(%v) instead",
			pool.MaxPoolSize)
		pool.BestPoolSize = pool.MaxPoolSize
	}

	pool.retryPool = make(chan int, pool.MaxPoolSize)
	pool.alivePool = make(chan interface{}, pool.MaxPoolSize)
	pool.swapPool = make(chan interface{}, pool.MaxPoolSize)

	for i := 0; i < pool.InitialPoolSize; i++ {
		if c, err := pool.Dial(ctx, pool.Address, pool.Port); err == nil {
			pool.alivePool <- c
		} else {
			pool.Logger.Errorf("initialize pool => Dial %v:%v failed,err:%v",
				pool.Address, pool.Port, err.Error())
			pool.retryPool <- 0
		}
	}

	go pool.start()

	return pool, nil
}

func (p *ExclusivePool) start() {
	go p.clean()
	go p.retryLoop()
	go p.keepAliveLoop()
}

func (p *ExclusivePool) Get() (connection interface{}, err error) {

	select {
	case connection = <-p.alivePool:
		{
			p.Logger.Infof("addr:%v => Get new connection from alive pool.", p.Address)
			atomic.AddInt32(&p.workConnCount, 1)
			return
		}
	case connection = <-p.swapPool:
		{
			p.Logger.Infof("addr:%v => Get new connection from swap pool.", p.Address)
			atomic.AddInt32(&p.workConnCount, 1)
			return
		}
	default:
		{
			p.Logger.Warnf("addr:%v => Get new connection from new create.", p.Address)

			if int(p.workConnCount)+len(p.retryPool)+len(p.alivePool)+len(p.swapPool) >= p.MaxPoolSize {
				p.Logger.Errorf("addr:%v => Pool Was Exhausted, detail: working: %v, alive: %v, retry: %v.",
					p.Address, p.workConnCount, len(p.alivePool), len(p.retryPool))
				return nil, ErrWrapper(ResourceExhausted, fmt.Sprintf("addr:%v", p.Address))
			}

			retry := 0
			for retry < p.DialRetryCount {
				if connection, err = p.Dial(context.TODO(), p.Address, p.Port); err != nil {
					p.Logger.Errorf("addr:%v => get conn => Dial %v:%v failed,err:%v",
						p.Address, p.Address, p.Port, err.Error())
					retry++
					continue
				} else {
					atomic.AddInt32(&p.workConnCount, 1)
					atomic.AddInt32(&p.newlyConnCount, 1)
					p.Logger.Errorf("addr:%v => get conn => Dial %v:%v successfully",
						p.Address, p.Address, p.Port)
					return
				}
			}

			if retry >= p.DialRetryCount {
				p.retryPool <- 0
				return nil, err
			}
		}
	}

	return nil, ErrWrapper(GetConnTimeout)
}

func (p *ExclusivePool) Put(connection interface{}) (err error) {

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

	return
}

func (p *ExclusivePool) Release() {
	p.isStopped = true

	for connection := range p.alivePool {
		if err := p.Close(context.TODO(), connection); err != nil {
			p.Logger.Infof("addr:%v => Release connection error:%v", p.Address, err)
		}
		atomic.SwapInt32(&p.workConnCount, p.workConnCount-1)
	}
}

func (p *ExclusivePool) retryLoop() {
	p.Logger.Infof("addr:%v => retry loop start.", p.Address)
	defer p.Logger.Infof("addr:%v => retry loop end.", p.Address)

	for {
		select {
		case <-time.After(p.DialRetryInterval):
			{
				max := len(p.retryPool)
				for i := 0; i < max; i++ {
					if connection, err := p.Dial(context.TODO(), p.Address, p.Port); err == nil {
						<-p.retryPool
						p.alivePool <- connection
						p.Logger.Infof("addr:%v => Retry Pool Success.", p.Address)
					} else {
						p.Logger.Errorf("addr:%v => Retry Pool Failed.", p.Address)
					}
				}

				if p.isStopped {
					break
				}
			}
		}
	}
}

func (p *ExclusivePool) keepAliveLoop() {

	p.Logger.Infof("addr:%v => keepAlive loop start.", p.Address)
	defer p.Logger.Infof("addr:%v => keepAlive loop end.", p.Address)

	for {
		select {
		case <-time.After(p.KeepAliveInterval):
			{
				if len(p.alivePool) > 0 {
					for connection := range p.alivePool {
						if err := p.KeepAlive(context.TODO(), connection); err == nil {
							p.swapPool <- connection
						} else {
							p.Logger.Errorf("addr:%v => Keepalive Pool Failed on %v\n",
								p.Address, fmt.Sprintf("%v:%v", p.Address, p.Port))
							p.retryPool <- 0
						}

						if len(p.alivePool) == 0 {
							break
						}
					}
				}

				if len(p.swapPool) > 0 {
					for connection := range p.swapPool {
						p.alivePool <- connection

						if len(p.swapPool) == 0 {
							break
						}
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
}

func (p *ExclusivePool) clean() {

	p.Logger.Infof("addr:%v => clean loop start.", p.Address)
	defer p.Logger.Infof("addr:%v => clean loop end.", p.Address)

	for {
		select {
		case <-time.After(p.CleanInterval):
			{
				if len(p.retryPool) > 0 {
					p.Logger.Infof("addr:%v the pool is retrying, skip.", p.Address)
					break
				}

				if atomic.LoadInt32(&p.newlyConnCount) > 0 {
					p.Logger.Infof("addr:%v the pool is at high load, skip.", p.Address)
					atomic.StoreInt32(&p.newlyConnCount, 0)
					break
				}

				for int(p.workConnCount)+len(p.alivePool)+len(p.swapPool) > p.BestPoolSize {
					select {
					case connection := <-p.alivePool:
						{
							p.Logger.Infof("addr:%v cleaning conn from alivePool.", p.Address)
							if err := p.Close(context.TODO(), connection); err != nil {
								p.Logger.Infof("addr:%v => cleaning connection error:%v", p.Address, err)
							}
						}
					case connection := <-p.swapPool:
						{
							p.Logger.Infof("addr:%v cleaning conn from swapPool.", p.Address)
							if err := p.Close(context.TODO(), connection); err != nil {
								p.Logger.Infof("addr:%v => cleaning connection error:%v", p.Address, err)
							}
						}
					}
				}
			}
		}
	}
}
