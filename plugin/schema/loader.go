package schema

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
)

var NewPluginObj NewPlugin

type Cfg struct {
	Viper  *viper.Viper
	Logger log.API
}

type NewPlugin = func(*Cfg) (Plugin, error)

type Loader interface {
	Init() error
	FInit() error
	Load(...NewPlugin) ([]Plugin, error)
}
