package rdb

import (
	"github.com/sqjian/go-kit/log"
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

	Logger log.API
}

func newDefaultMeta() *DbMeta {
	return &DbMeta{
		Logger: log.DummyLogger,
	}
}

func newMeta(opts ...MetaOption) *DbMeta {

	meta := newDefaultMeta()

	for _, opt := range opts {
		opt.apply(meta)
	}
	return meta
}
