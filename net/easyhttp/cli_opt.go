package easyhttp

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"net/http"
)

type cliCfg struct {
	logId string

	retry   int
	trace   bool
	logger  log.API
	client  *http.Client
	context context.Context

	body   []byte
	query  map[string]string
	header map[string]string
}

type CliOption interface {
	apply(*cliCfg)
}

type cliOptionFunc func(*cliCfg)

func (f cliOptionFunc) apply(options *cliCfg) {
	f(options)
}

func WithCliHeader(header map[string]string) CliOption {
	return cliOptionFunc(func(options *cliCfg) {
		options.header = header
	})
}

func WithCliQuery(query map[string]string) CliOption {
	return cliOptionFunc(func(options *cliCfg) {
		options.query = query
	})
}

func WithCliBody(body []byte) CliOption {
	return cliOptionFunc(func(options *cliCfg) {
		options.body = body
	})
}

func WithClient(client *http.Client) CliOption {
	return cliOptionFunc(func(options *cliCfg) {
		options.client = client
	})
}

func WithCliUniqueId(uniqueId string) CliOption {
	return cliOptionFunc(func(options *cliCfg) {
		options.logId = uniqueId
	})
}

func WithCliRetry(retry int) CliOption {
	return cliOptionFunc(func(options *cliCfg) {
		options.retry = retry
	})
}

func WithCliLogger(logger log.API) CliOption {
	return cliOptionFunc(func(options *cliCfg) {
		options.logger = logger
	})
}

func WithCliTrace(trace bool) CliOption {
	return cliOptionFunc(func(options *cliCfg) {
		options.trace = trace
	})
}
