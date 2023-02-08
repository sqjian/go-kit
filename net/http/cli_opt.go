package http

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"net/http"
)

type cliConfig struct {
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

type CliOptionFunc func(*cliConfig)

func WithCliHeader(header map[string]string) CliOptionFunc {
	return func(options *cliConfig) {
		options.header = header
	}
}

func WithCliQuery(query map[string]string) CliOptionFunc {
	return func(options *cliConfig) {
		options.query = query
	}
}

func WithCliBody(body []byte) CliOptionFunc {
	return func(options *cliConfig) {
		options.body = body
	}
}

func WithClient(client *http.Client) CliOptionFunc {
	return func(options *cliConfig) {
		options.client = client
	}
}

func WithCliUniqueId(uniqueId string) CliOptionFunc {
	return func(options *cliConfig) {
		options.logId = uniqueId
	}
}

func WithCliRetry(retry int) CliOptionFunc {
	return func(options *cliConfig) {
		options.retry = retry
	}
}

func WithCliLogger(logger log.API) CliOptionFunc {
	return func(options *cliConfig) {
		options.logger = logger
	}
}

func WithCliTrace(trace bool) CliOptionFunc {
	return func(options *cliConfig) {
		options.trace = trace
	}
}
