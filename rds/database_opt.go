package rds

import (
	"github.com/sqjian/go-kit/log"
	"time"
)

type ConfigOptionFunc func(*Config)

func WithLogger(logger log.Log) ConfigOptionFunc {
	return func(m *Config) {
		m.Logger = logger
	}
}

func WithUserName(UserName string) ConfigOptionFunc {
	return func(m *Config) {
		m.UserName = UserName
	}
}

func WithPassWord(PassWord string) ConfigOptionFunc {
	return func(m *Config) {
		m.PassWord = PassWord
	}
}

func WithIp(ip string) ConfigOptionFunc {
	return func(m *Config) {
		m.IP = ip
	}
}

func WithPort(port string) ConfigOptionFunc {
	return func(m *Config) {
		m.Port = port
	}
}

func WithDbName(dbName string) ConfigOptionFunc {
	return func(m *Config) {
		m.DbName = dbName
	}
}

func WithMaxIdleConns(MaxIdleConns int) ConfigOptionFunc {
	return func(m *Config) {
		m.MaxIdleConns = MaxIdleConns
	}
}

func WithMaxLifeTime(MaxLifeTime time.Duration) ConfigOptionFunc {
	return func(m *Config) {
		m.MaxLifeTime = MaxLifeTime
	}
}
