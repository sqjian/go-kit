package rdb

import "time"

type Meta struct {
	IP       string
	Port     string
	DbName   string
	UserName string
	PassWord string

	MaxIdleConns int
	MaxLifeTime  time.Duration
}

func newDefaultMeta() *Meta {
	return &Meta{}
}

func NewMeta(opts ...Option) {

	meta := newDefaultMeta()

	for _, opt := range opts {
		opt.apply(meta)
	}

	panic("IMPL")
}
