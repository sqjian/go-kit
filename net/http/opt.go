package http

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"net/http"
)

type Config struct {
	logId string

	retry   int
	logger  log.Logger
	client  *http.Client
	context context.Context

	body   []byte
	query  map[string]string
	header map[string]string
}

type Option interface {
	apply(*Config)
}

type optionFunc func(*Config)

func (f optionFunc) apply(options *Config) {
	f(options)
}

func WithHeader(header map[string]string) Option {
	return optionFunc(func(options *Config) {
		options.header = header
	})
}

func WithQuery(query map[string]string) Option {
	return optionFunc(func(options *Config) {
		options.query = query
	})
}

func WithBody(body []byte) Option {
	return optionFunc(func(options *Config) {
		options.body = body
	})
}

func WithClient(client *http.Client) Option {
	return optionFunc(func(options *Config) {
		options.client = client
	})
}

func WithUniqueId(uniqueId string) Option {
	return optionFunc(func(options *Config) {
		options.logId = uniqueId
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
