package rdb

import (
	"github.com/sqjian/go-kit/log"
	"time"
)

type Meta struct {
	IP       string
	Port     string
	DbName   string
	UserName string
	PassWord string

	MaxIdleConns int
	MaxLifeTime  time.Duration

	Logger log.API
}

func newDefaultMeta() *Meta {
	return &Meta{
		Logger: log.DummyLogger,
	}
}

func NewMeta(opts ...Option) {

	meta := newDefaultMeta()

	for _, opt := range opts {
		opt.apply(meta)
	}

	panic("IMPL")
}