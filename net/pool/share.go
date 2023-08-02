package pool

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"sync"
	"sync/atomic"
	"time"
)

func newSharePool(cfg *Config) *SharePool {
	return &SharePool{
		Address:           cfg.Address,
		Port:              cfg.Port,
		Dial:              cfg.Dial,
		Close:             cfg.Close,
		KeepAlive:         cfg.KeepAlive,
		InitialPoolSize:   cfg.InitialPoolSize,
		DialRetryCount:    cfg.DialRetryCount,
		KeepAliveInterval: cfg.KeepAliveInterval,
		CleanInterval:     cfg.CleanInterval,
		DialRetryInterval: cfg.DialRetryInterval,
		CreateNewInterval: cfg.CreateNewInterval,
		Logger:            cfg.Logger,
	}
}

type SharePool struct {
	Address           string
	Port              string
	Dial              func(ctx context.Context, address, port string) (connection any, err error)
	Close             func(ctx context.Context, connection any) (err error)
	KeepAlive         func(ctx context.Context, connection any) (err error)
	InitialPoolSize   int
	DialRetryCount    int
	KeepAliveInterval time.Duration
	CleanInterval     time.Duration
	DialRetryInterval time.Duration
	CreateNewInterval time.Duration
	Logger            log.Log

	alivePool       []any
	alivePoolOffset uint64

	retryPool chan int
	sync      sync.RWMutex
	isStopped bool
}

func NewSharePool(ctx context.Context, cfg *Config) (*SharePool, error) {

	pool := newSharePool(cfg)

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

	for i := 0; i < pool.InitialPoolSize; i++ {
		if c, err := pool.Dial(ctx, pool.Address, pool.Port); err == nil {
			pool.alivePool = append(pool.alivePool, c)
		} else {
			pool.Logger.Errorf("initialize pool => Dial %v:%v failed,err:%v",
				pool.Address, pool.Port, err.Error())
			pool.retryPool <- 0
		}
	}

	go pool.start()

	return pool, nil
}

func (s *SharePool) start() {
	go s.retryLoop()
	go s.keepAliveLoop()
}

func (s *SharePool) keepAliveLoop() {

	s.Logger.Infof("addr:%v => keepAlive loop start.", s.Address)
	defer s.Logger.Infof("addr:%v => keepAlive loop end.", s.Address)

	remove := func(array []any, offset int) []any {
		for j := offset; j < len(array)-1; j++ {
			array[j] = array[j+1]
		}
		return array[:len(array)-1]
	}
	pickOut := func(array []any) int {
		for offset, connection := range s.alivePool {
			if err := s.KeepAlive(context.TODO(), connection); err != nil {
				s.Logger.Errorf("addr:%v => Keepalive Pool Failed on %v\n",
					s.Address, fmt.Sprintf("%v:%v", s.Address, s.Port))
				return offset
			}
		}
		return 0
	}
	for {
		select {
		case <-time.After(s.KeepAliveInterval):
			{
				for {
					if len(s.alivePool) > 0 {
						s.sync.RLock()
						offset := pickOut(s.alivePool)
						s.sync.RUnlock()
						if offset == 0 {
							break
						}

						s.retryPool <- 0
						s.sync.Lock()
						s.alivePool = remove(s.alivePool, offset)
						s.sync.Unlock()
					}
				}
			}
		}

		if s.isStopped {
			for connection := range s.alivePool {
				err := s.Close(context.TODO(), connection)
				if err != nil {
					return
				}
			}
			break
		}
	}
}
func (s *SharePool) retryLoop() {
	s.Logger.Infof("addr:%v => retry loop start.", s.Address)
	defer s.Logger.Infof("addr:%v => retry loop end.", s.Address)

	for {
		select {
		case <-time.After(s.DialRetryInterval):
			{
				max := len(s.retryPool)
				for i := 0; i < max; i++ {
					if connection, err := s.Dial(context.TODO(), s.Address, s.Port); err == nil {
						<-s.retryPool
						s.sync.Lock()
						s.alivePool = append(s.alivePool, connection)
						s.sync.Unlock()
						s.Logger.Infof("addr:%v => Retry Pool Success.", s.Address)
					} else {
						s.Logger.Errorf("addr:%v => Retry Pool Failed.", s.Address)
					}
				}

				if s.isStopped {
					break
				}
			}
		}
	}
}
func (s *SharePool) Get() (connection any, err error) {
	s.sync.RLock()
	defer s.sync.RUnlock()

	if len(s.alivePool) == 0 {
		s.Logger.Errorf("addr:%v => Pool Was Exhausted, detail: alive: %v, retry: %v.",
			s.Address, len(s.alivePool), len(s.retryPool))

		return nil, ErrWrapper(ResourceExhausted, fmt.Sprintf("addr:%v", s.Address))
	}

	return s.alivePool[atomic.AddUint64(&s.alivePoolOffset, 1)%uint64(len(s.alivePool))], nil
}

func (s *SharePool) Put(_ any) (err error) {
	return nil
}

func (s *SharePool) Release() {
	s.sync.Lock()
	s.isStopped = true

	for _, connection := range s.alivePool {
		if err := s.Close(context.TODO(), connection); err != nil {
			s.Logger.Infof("addr:%v => Release connection error:%v", s.Address, err)
		}
	}

	s.sync.Unlock()
}
