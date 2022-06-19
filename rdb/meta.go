package rdb

import (
	"github.com/sqjian/go-kit/easylog"
	"time"
)

type DbMeta struct {
	IP       string
	Port     string
	DbName   string
	UserName string
	PassWord string

	MaxIdleConns int
	MaxLifeTime  time.Duration

	Logger easylog.API
}

func newDefaultMeta() *DbMeta {
	return &DbMeta{
		Logger: easylog.DummyLogger,
	}
}

func newMeta(opts ...MetaOption) *DbMeta {

	meta := newDefaultMeta()

	for _, opt := range opts {
		opt.apply(meta)
	}
	return meta
}
