package http

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"net/http"
)

type clientConfig struct {
	log.Log

	retry   int
	trace   bool
	client  *http.Client
	context context.Context

	body   []byte
	query  map[string]string
	header map[string]string
}

type CliOptionFunc func(*clientConfig)

func WithCliHeader(header map[string]string) CliOptionFunc {
	return func(options *clientConfig) {
		options.header = header
	}
}

func WithCliQuery(query map[string]string) CliOptionFunc {
	return func(options *clientConfig) {
		options.query = query
	}
}

func WithCliBody(body []byte) CliOptionFunc {
	return func(options *clientConfig) {
		options.body = body
	}
}

func WithClient(client *http.Client) CliOptionFunc {
	return func(options *clientConfig) {
		options.client = client
	}
}

func WithCliRetry(retry int) CliOptionFunc {
	return func(options *clientConfig) {
		options.retry = retry
	}
}

func WithCliLogger(logger log.Log) CliOptionFunc {
	return func(options *clientConfig) {
		if logger == nil {
			panic("please check if the log is initialized")
		}
		options.Log = logger
	}
}

func WithCliTrace(trace bool) CliOptionFunc {
	return func(options *clientConfig) {
		options.trace = trace
	}
}
