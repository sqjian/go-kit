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

func newDefaultCfg() *Config {
	return &Config{
		PoolType:          Exclusive,
		Logger:            func() log.Log { inst, _ := log.NewLogger(log.WithLevel("dummy")); return inst }(),
		KeepAliveInterval: DefaultKeepAliveInterval,
		CreateNewInterval: DefaultCreateNewInterval,
		DialRetryCount:    DefaultDialRetryCount,
		DialRetryInterval: DefaultRetryInterval,
		CleanInterval:     DefaultCleanInterval,
	}
}

type Config struct {
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
	Logger            log.Log
}

type OptionFunc func(*Config)

func WithType(poolType Type) OptionFunc {
	return func(options *Config) {
		options.PoolType = poolType
	}
}

func WithAddress(Address string) OptionFunc {
	return func(options *Config) {
		options.Address = Address
	}
}

func WithPort(Port string) OptionFunc {
	return func(options *Config) {
		options.Port = Port
	}
}

func WithDial(Dial func(ctx context.Context, address, port string) (connection interface{}, err error)) OptionFunc {
	return func(options *Config) {
		options.Dial = Dial
	}
}

func WithClose(Close func(ctx context.Context, connection interface{}) (err error)) OptionFunc {
	return func(options *Config) {
		options.Close = Close
	}
}

func WithKeepAlive(KeepAlive func(ctx context.Context, connection interface{}) (err error)) OptionFunc {
	return func(options *Config) {
		options.KeepAlive = KeepAlive
	}
}

func WithInitialPoolSize(InitialPoolSize int) OptionFunc {
	return func(options *Config) {
		options.InitialPoolSize = InitialPoolSize
	}
}

func WithBestPoolSize(BestPoolSize int) OptionFunc {
	return func(options *Config) {
		options.BestPoolSize = BestPoolSize
	}
}

func WithMaxPoolSize(MaxPoolSize int) OptionFunc {
	return func(options *Config) {
		options.MaxPoolSize = MaxPoolSize
	}
}

func WithDialRetryCount(DialRetryCount int) OptionFunc {
	return func(options *Config) {
		options.DialRetryCount = DialRetryCount
	}
}

func WithKeepAliveInterval(KeepAliveInterval time.Duration) OptionFunc {
	return func(options *Config) {
		options.KeepAliveInterval = KeepAliveInterval
	}
}

func WithCleanInterval(CleanInterval time.Duration) OptionFunc {
	return func(options *Config) {
		options.CleanInterval = CleanInterval
	}
}

func WithDialRetryInterval(DialRetryInterval time.Duration) OptionFunc {
	return func(options *Config) {
		options.DialRetryInterval = DialRetryInterval
	}
}

func WithCreateNewInterval(CreateNewInterval time.Duration) OptionFunc {
	return func(options *Config) {
		options.CreateNewInterval = CreateNewInterval
	}
}

func WithLogger(Logger log.Log) OptionFunc {
	return func(options *Config) {
		options.Logger = Logger
	}
}
