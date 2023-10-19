package rds

import (
	"github.com/sqjian/go-kit/log"
	"time"
)

type config struct {
	IP       string
	Port     string
	DbName   string
	UserName string
	PassWord string

	MaxIdleConns int
	MaxLifeTime  time.Duration

	placeHolder string

	Logger log.Log
}

func newDefaultConfig() *config {
	return func() *config {
		return &config{
			Logger: func() log.Log { inst, _ := log.NewLogger(log.WithLevel("dummy")); return inst }(),
		}
	}()
}

type ConfigOptionFunc func(*config)

func WithLogger(logger log.Log) ConfigOptionFunc {
	return func(m *config) {
		m.Logger = logger
	}
}

func WithUserName(UserName string) ConfigOptionFunc {
	return func(m *config) {
		m.UserName = UserName
	}
}

func WithPassWord(PassWord string) ConfigOptionFunc {
	return func(m *config) {
		m.PassWord = PassWord
	}
}

func WithIp(ip string) ConfigOptionFunc {
	return func(m *config) {
		m.IP = ip
	}
}

func WithPort(port string) ConfigOptionFunc {
	return func(m *config) {
		m.Port = port
	}
}

func WithDbName(dbName string) ConfigOptionFunc {
	return func(m *config) {
		m.DbName = dbName
	}
}

func WithMaxIdleConns(MaxIdleConns int) ConfigOptionFunc {
	return func(m *config) {
		m.MaxIdleConns = MaxIdleConns
	}
}

func WithMaxLifeTime(MaxLifeTime time.Duration) ConfigOptionFunc {
	return func(m *config) {
		m.MaxLifeTime = MaxLifeTime
	}
}
