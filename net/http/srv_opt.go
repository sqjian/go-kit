package http

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"time"
)

type serverConfig struct {
	limit          int
	MaxHeaderBytes int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	logger     log.Log
	gracefully time.Duration
	context    context.Context
}

type ServerOptionFunc func(*serverConfig)

func WithSrvLogger(logger log.Log) ServerOptionFunc {
	return func(options *serverConfig) {
		options.logger = logger
	}
}
