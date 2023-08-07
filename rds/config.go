package rds

import (
	"github.com/sqjian/go-kit/log"
	"time"
)

type Config struct {
	IP       string
	Port     string
	DbName   string
	UserName string
	PassWord string

	MaxIdleConns int
	MaxLifeTime  time.Duration

	Logger log.Log
}

func newDefaultConfig() *Config {
	return &Config{
		Logger: func() log.Log { inst, _ := log.NewLogger(log.WithLevel("dummy")); return inst }(),
	}
}

func newConfig(opts ...ConfigOptionFunc) *Config {

	config := newDefaultConfig()

	for _, opt := range opts {
		opt(config)
	}
	return config
}
