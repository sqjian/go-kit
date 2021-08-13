package connection

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"time"
)

type option interface {
	apply(*ClientPool)
}

type OptionFunc func(*ClientPool)

func (f OptionFunc) apply(options *ClientPool) {
	f(options)
}

func WithAddress(Address string) option {
	return OptionFunc(func(options *ClientPool) {
		options.Address = Address
	})
}

func WithPort(Port string) option {
	return OptionFunc(func(options *ClientPool) {
		options.Port = Port
	})
}

func WithDial(Dial func(ctx context.Context, address, port string) (connection interface{}, err error)) option {
	return OptionFunc(func(options *ClientPool) {
		options.Dial = Dial
	})
}

func WithClose(Close func(ctx context.Context, connection interface{}) (err error)) option {
	return OptionFunc(func(options *ClientPool) {
		options.Close = Close
	})
}

func WithKeepAlive(KeepAlive func(ctx context.Context, connection interface{}) (err error)) option {
	return OptionFunc(func(options *ClientPool) {
		options.KeepAlive = KeepAlive
	})
}

func WithInitialPoolSize(InitialPoolSize int) option {
	return OptionFunc(func(options *ClientPool) {
		options.InitialPoolSize = InitialPoolSize
	})
}

func WithMaxPoolSize(MaxPoolSize int) option {
	return OptionFunc(func(options *ClientPool) {
		options.MaxPoolSize = MaxPoolSize
	})
}

func WithDialRetryCount(DialRetryCount int) option {
	return OptionFunc(func(options *ClientPool) {
		options.DialRetryCount = DialRetryCount
	})
}

func WithKeepAliveInterval(KeepAliveInterval time.Duration) option {
	return OptionFunc(func(options *ClientPool) {
		options.KeepAliveInterval = KeepAliveInterval
	})
}

func WithDialRetryInterval(DialRetryInterval time.Duration) option {
	return OptionFunc(func(options *ClientPool) {
		options.DialRetryInterval = DialRetryInterval
	})
}

func WithCreateNewInterval(CreateNewInterval time.Duration) option {
	return OptionFunc(func(options *ClientPool) {
		options.CreateNewInterval = CreateNewInterval
	})
}

func WithLogger(Logger log.Logger) option {
	return OptionFunc(func(options *ClientPool) {
		options.Logger = Logger
	})
}
