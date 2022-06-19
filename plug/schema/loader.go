package schema

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/easylog"
)

var NewPlugObj NewPlug

type Cfg struct {
	Viper  *viper.Viper
	Logger easylog.API
}

type NewPlug = func(func(*Cfg)) (Plug, error)

type Loader interface {
	Init() error
	FInit() error
	Load(...NewPlug) ([]Plug, error)
}
