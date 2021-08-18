package retry

import (
	"context"
	"time"
)

type IfFunc func(error) bool

type OnRetryFunc func(n uint, err error)

type Config struct {
	attempts  uint
	delayTime time.Duration
	onRetry   OnRetryFunc
	retryIf   IfFunc
	context   context.Context
}

type Option func(*Config)

func WithAttempts(attempts uint) Option {
	return func(c *Config) {
		c.attempts = attempts
	}
}

func WithDelay(delay time.Duration) Option {
	return func(c *Config) {
		c.delayTime = delay
	}
}

func WithOnRetry(onRetry OnRetryFunc) Option {
	return func(c *Config) {
		c.onRetry = onRetry
	}
}

func WithIf(retryIf IfFunc) Option {
	return func(c *Config) {
		c.retryIf = retryIf
	}
}

func WithContext(ctx context.Context) Option {
	return func(c *Config) {
		c.context = ctx
	}
}
