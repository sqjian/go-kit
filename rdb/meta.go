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
		Logger: log.DummyLogger{},
	}
}

func newMeta(opts ...MetaOptionFunc) *Meta {

	meta := newDefaultMeta()

	for _, opt := range opts {
		opt(meta)
	}
	return meta
}
