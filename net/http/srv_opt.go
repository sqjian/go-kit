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

type SrvOption interface {
	apply(*srvCfg)
}

type srvOptionFunc func(*srvCfg)

func (f srvOptionFunc) apply(options *srvCfg) {
	f(options)
}

func WithSrvLogger(logger log.API) SrvOption {
	return srvOptionFunc(func(options *srvCfg) {
		options.logger = logger
	})
}
