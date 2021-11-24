package pool

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"time"
)

var (
	DefaultDialRetryCount    = 3
	DefaultRetryInterval     = time.Second * 10
	DefaultKeepAliveInterval = time.Second * 3
	DefaultCreateNewInterval = time.Second * 1
	DefaultCleanInterval     = time.Second * 60
)

func newDefaultCfg() *Cfg {
	return &Cfg{
		PoolType:          Exclusive,
		Logger:            log.DummyLogger,
		KeepAliveInterval: DefaultKeepAliveInterval,
		CreateNewInterval: DefaultCreateNewInterval,
		DialRetryCount:    DefaultDialRetryCount,
		DialRetryInterval: DefaultRetryInterval,
		CleanInterval:     DefaultCleanInterval,
	}
}

type Cfg struct {
	PoolType          Type
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
	CreateNewInterval time.Duration
	Logger            log.API
}

type Option interface {
	apply(*Cfg)
}

type OptionFunc func(*Cfg)

func (f OptionFunc) apply(options *Cfg) {
	f(options)
}

func WithType(poolType Type) Option {
	return OptionFunc(func(options *Cfg) {
		options.PoolType = poolType
	})
}

func WithAddress(Address string) Option {
	return OptionFunc(func(options *Cfg) {
		options.Address = Address
	})
}

func WithPort(Port string) Option {
	return OptionFunc(func(options *Cfg) {
		options.Port = Port
	})
}

func WithDial(Dial func(ctx context.Context, address, port string) (connection interface{}, err error)) Option {
	return OptionFunc(func(options *Cfg) {
		options.Dial = Dial
	})
}

func WithClose(Close func(ctx context.Context, connection interface{}) (err error)) Option {
	return OptionFunc(func(options *Cfg) {
		options.Close = Close
	})
}

func WithKeepAlive(KeepAlive func(ctx context.Context, connection interface{}) (err error)) Option {
	return OptionFunc(func(options *Cfg) {
		options.KeepAlive = KeepAlive
	})
}

func WithInitialPoolSize(InitialPoolSize int) Option {
	return OptionFunc(func(options *Cfg) {
		options.InitialPoolSize = InitialPoolSize
	})
}

func WithBestPoolSize(BestPoolSize int) Option {
	return OptionFunc(func(options *Cfg) {
		options.BestPoolSize = BestPoolSize
	})
}

func WithMaxPoolSize(MaxPoolSize int) Option {
	return OptionFunc(func(options *Cfg) {
		options.MaxPoolSize = MaxPoolSize
	})
}

func WithDialRetryCount(DialRetryCount int) Option {
	return OptionFunc(func(options *Cfg) {
		options.DialRetryCount = DialRetryCount
	})
}

func WithKeepAliveInterval(KeepAliveInterval time.Duration) Option {
	return OptionFunc(func(options *Cfg) {
		options.KeepAliveInterval = KeepAliveInterval
	})
}

func WithCleanInterval(CleanInterval time.Duration) Option {
	return OptionFunc(func(options *Cfg) {
		options.CleanInterval = CleanInterval
	})
}

func WithDialRetryInterval(DialRetryInterval time.Duration) Option {
	return OptionFunc(func(options *Cfg) {
		options.DialRetryInterval = DialRetryInterval
	})
}

func WithCreateNewInterval(CreateNewInterval time.Duration) Option {
	return OptionFunc(func(options *Cfg) {
		options.CreateNewInterval = CreateNewInterval
	})
}

func WithLogger(Logger log.API) Option {
	return OptionFunc(func(options *Cfg) {
		options.Logger = Logger
	})
}
