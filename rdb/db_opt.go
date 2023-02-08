package rdb

import (
	"github.com/sqjian/go-kit/log"
	"time"
)

type MetaOptionFunc func(*Meta)

func WithLogger(logger log.API) MetaOptionFunc {
	return func(m *Meta) {
		m.Logger = logger
	}
}

func WithUserName(UserName string) MetaOptionFunc {
	return func(m *Meta) {
		m.UserName = UserName
	}
}

func WithPassWord(PassWord string) MetaOptionFunc {
	return func(m *Meta) {
		m.PassWord = PassWord
	}
}

func WithIp(ip string) MetaOptionFunc {
	return func(m *Meta) {
		m.IP = ip
	}
}

func WithPort(port string) MetaOptionFunc {
	return func(m *Meta) {
		m.Port = port
	}
}

func WithDbName(dbName string) MetaOptionFunc {
	return func(m *Meta) {
		m.DbName = dbName
	}
}

func WithMaxIdleConns(MaxIdleConns int) MetaOptionFunc {
	return func(m *Meta) {
		m.MaxIdleConns = MaxIdleConns
	}
}

func WithMaxLifeTime(MaxLifeTime time.Duration) MetaOptionFunc {
	return func(m *Meta) {
		m.MaxLifeTime = MaxLifeTime
	}
}
