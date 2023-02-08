package http

import (
	"context"
	"github.com/sqjian/go-kit/log"
	"time"
)

type srvCfg struct {
	logId string

	limit          int
	MaxHeaderBytes int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	logger     log.API
	gracefully time.Duration
	context    context.Context
}

type SrvOptionFunc func(*srvCfg)

func WithSrvLogger(logger log.API) SrvOptionFunc {
	return func(options *srvCfg) {
		options.logger = logger
	}
}
