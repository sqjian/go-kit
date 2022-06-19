package retry

import (
	"context"
	"github.com/sqjian/go-kit/easylog"
	"time"
)

type IfFunc func(error) bool

type OnRetryFunc func(n uint, err error)

type DelayTypeFunc func(n uint, err error, config *Config) time.Duration

type Config struct {
	attempts  uint
	retryIf   IfFunc
	onRetry   OnRetryFunc
	delayType DelayTypeFunc
	context   context.Context
	logger    easylog.API
}

type Option func(*Config)

func WithAttempts(attempts uint) Option {
	return func(c *Config) {
		c.attempts = attempts
	}
}

func WithDelayFn(delayTypeFn DelayTypeFunc) Option {
	return func(c *Config) {
		c.delayType = delayTypeFn
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
