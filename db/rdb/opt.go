package rdb

import (
	"github.com/sqjian/go-kit/log"
	"time"
)

type Option interface {
	apply(*Meta)
}

type optionFunc func(*Meta)

func (f optionFunc) apply(log *Meta) {
	f(log)
}

func WithLogger(logger log.API) Option {
	return optionFunc(func(m *Meta) {
		m.Logger = logger
	})
}

func WithUserName(UserName string) Option {
	return optionFunc(func(m *Meta) {
		m.UserName = UserName
	})
}

func WithPassWord(PassWord string) Option {
	return optionFunc(func(m *Meta) {
		m.PassWord = PassWord
	})
}

func WithIp(ip string) Option {
	return optionFunc(func(m *Meta) {
		m.IP = ip
	})
}

func WithPort(port string) Option {
	return optionFunc(func(m *Meta) {
		m.Port = port
	})
}

func WithDbName(dbName string) Option {
	return optionFunc(func(m *Meta) {
		m.DbName = dbName
	})
}

func WithMaxIdleConns(MaxIdleConns int) Option {
	return optionFunc(func(m *Meta) {
		m.MaxIdleConns = MaxIdleConns
	})
}

func WithMaxLifeTime(MaxLifeTime time.Duration) Option {
	return optionFunc(func(m *Meta) {
		m.MaxLifeTime = MaxLifeTime
	})
}
