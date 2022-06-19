package easyhttp

import (
	"context"
	"github.com/sqjian/go-kit/easylog"
	"time"
)

type srvCfg struct {
	logId string

	limit          int
	MaxHeaderBytes int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	logger     easylog.API
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

func WithSrvLogger(logger easylog.API) SrvOption {
	return srvOptionFunc(func(options *srvCfg) {
		options.logger = logger
	})
}
