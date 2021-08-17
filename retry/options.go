package retry

import (
	"context"
	"math"
	"math/rand"
	"time"
)

type IfFunc func(error) bool

type OnRetryFunc func(n uint, err error)

type DelayTypeFunc func(n uint, err error, config *Config) time.Duration

type Config struct {
	attempts      uint
	delay         time.Duration
	maxDelay      time.Duration
	maxJitter     time.Duration
	onRetry   OnRetryFunc
	retryIf   IfFunc
	delayType DelayTypeFunc
	lastErrorOnly bool
	context       context.Context

	maxBackOffN uint
}

type Option func(*Config)

func LastErrorOnly(lastErrorOnly bool) Option {
	return func(c *Config) {
		c.lastErrorOnly = lastErrorOnly
	}
}

func Attempts(attempts uint) Option {
	return func(c *Config) {
		c.attempts = attempts
	}
}

func Delay(delay time.Duration) Option {
	return func(c *Config) {
		c.delay = delay
	}
}

func MaxDelay(maxDelay time.Duration) Option {
	return func(c *Config) {
		c.maxDelay = maxDelay
	}
}

func MaxJitter(maxJitter time.Duration) Option {
	return func(c *Config) {
		c.maxJitter = maxJitter
	}
}

func DelayType(delayType DelayTypeFunc) Option {
	return func(c *Config) {
		c.delayType = delayType
	}
}

func BackOffDelay(n uint, _ error, config *Config) time.Duration {
	const max uint = 62

	if config.maxBackOffN == 0 {
		if config.delay <= 0 {
			config.delay = 1
		}

		config.maxBackOffN = max - uint(math.Floor(math.Log2(float64(config.delay))))
	}

	if n > config.maxBackOffN {
		n = config.maxBackOffN
	}

	return config.delay << n
}

func FixedDelay(_ uint, _ error, config *Config) time.Duration {
	return config.delay
}

func RandomDelay(_ uint, _ error, config *Config) time.Duration {
	return time.Duration(rand.Int63n(int64(config.maxJitter)))
}

func CombineDelay(delays ...DelayTypeFunc) DelayTypeFunc {
	const maxInt64 = uint64(math.MaxInt64)

	return func(n uint, err error, config *Config) time.Duration {
		var total uint64
		for _, delay := range delays {
			total += uint64(delay(n, err, config))
			if total > maxInt64 {
				total = maxInt64
			}
		}

		return time.Duration(total)
	}
}

func OnRetry(onRetry OnRetryFunc) Option {
	return func(c *Config) {
		c.onRetry = onRetry
	}
}

func If(retryIf IfFunc) Option {
	return func(c *Config) {
		c.retryIf = retryIf
	}
}

func Context(ctx context.Context) Option {
	return func(c *Config) {
		c.context = ctx
	}
}
