package schema

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
)

var NewPlugObj NewPlug

type Cfg struct {
	Viper  *viper.Viper
	Logger log.Log
}

type NewPlug = func(func(*Cfg)) (Plug, error)

type Loader interface {
	Init() error
	FInit() error
	Load(...NewPlug) ([]Plug, error)
}
