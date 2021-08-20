package http

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"net/http"
)

type Config struct {
	retry   int
	client  *http.Client
	context context.Context
	logger  log.Logger
}

type Option interface {
	apply(*Config)
}

type optionFunc func(*Config)

func (f optionFunc) apply(options *Config) {
	f(options)
}

func WithClient(client *http.Client) Option {
	return optionFunc(func(options *Config) {
		options.client = client
	})
}

func WithContext(ctx context.Context) Option {
	return optionFunc(func(options *Config) {
		options.context = ctx
	})
}

func WithRetry(retry int) Option {
	return optionFunc(func(options *Config) {
		options.retry = retry
	})
}

func WithLogger(Logger log.Logger) Option {
	return optionFunc(func(options *Config) {
		options.logger = Logger
	})
}
