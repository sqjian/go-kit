package http

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/retry"
	"net/http"
)

type clientConfig struct {
	log.Log

	retry     int
	delayType retry.DelayTypeFunc

	trace   bool
	client  *http.Client
	context context.Context

	body   []byte
	query  map[string]string
	header map[string]string
}

type ClientOptionFunc func(*clientConfig)

func WithClientHeader(header map[string]string) ClientOptionFunc {
	return func(options *clientConfig) {
		options.header = header
	}
}

func WithClientQuery(query map[string]string) ClientOptionFunc {
	return func(options *clientConfig) {
		options.query = query
	}
}

func WithClientBody(body []byte) ClientOptionFunc {
	return func(options *clientConfig) {
		options.body = body
	}
}

func WithClient(client *http.Client) ClientOptionFunc {
	return func(options *clientConfig) {
		options.client = client
	}
}

func WithClientRetry(retry int) ClientOptionFunc {
	return func(options *clientConfig) {
		options.retry = retry
	}
}

func WithClientDelayFn(delayTypeFn retry.DelayTypeFunc) ClientOptionFunc {
	return func(options *clientConfig) {
		options.delayType = delayTypeFn
	}
}

func WithClientLogger(logger log.Log) ClientOptionFunc {
	return func(options *clientConfig) {
		if logger == nil {
			panic("please check if the log is initialized")
		}
		options.Log = logger
	}
}

func WithClientTrace(trace bool) ClientOptionFunc {
	return func(options *clientConfig) {
		options.trace = trace
	}
}
