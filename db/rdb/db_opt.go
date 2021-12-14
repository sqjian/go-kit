package rdb

import (
	"github.com/sqjian/go-kit/log"
	"time"
)

type MetaOption interface {
	apply(*DbMeta)
}

type metaOptionFunc func(*DbMeta)

func (f metaOptionFunc) apply(log *DbMeta) {
	f(log)
}

func WithLogger(logger log.API) MetaOption {
	return metaOptionFunc(func(m *DbMeta) {
		m.Logger = logger
	})
}

func WithUserName(UserName string) MetaOption {
	return metaOptionFunc(func(m *DbMeta) {
		m.UserName = UserName
	})
}

func WithPassWord(PassWord string) MetaOption {
	return metaOptionFunc(func(m *DbMeta) {
		m.PassWord = PassWord
	})
}

func WithIp(ip string) MetaOption {
	return metaOptionFunc(func(m *DbMeta) {
		m.IP = ip
	})
}

func WithPort(port string) MetaOption {
	return metaOptionFunc(func(m *DbMeta) {
		m.Port = port
	})
}

func WithDbName(dbName string) MetaOption {
	return metaOptionFunc(func(m *DbMeta) {
		m.DbName = dbName
	})
}

func WithMaxIdleConns(MaxIdleConns int) MetaOption {
	return metaOptionFunc(func(m *DbMeta) {
		m.MaxIdleConns = MaxIdleConns
	})
}

func WithMaxLifeTime(MaxLifeTime time.Duration) MetaOption {
	return metaOptionFunc(func(m *DbMeta) {
		m.MaxLifeTime = MaxLifeTime
	})
}
