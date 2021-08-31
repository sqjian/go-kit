package loader

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/schema"
	"sync"
)

type Cfg struct {
	Viper  *viper.Viper
	Logger log.Logger
}

type Loader struct {
	sync.RWMutex
	viper  *viper.Viper
	logger log.Logger
}

func NewLoader(cfg *Cfg) schema.Loader {
	return &Loader{viper: cfg.Viper, logger: cfg.Logger}
}

func (l *Loader) Init() error {
	return nil
}

func (l *Loader) FInit() error {
	return nil
}

func (l *Loader) Load(pluginGenerators ...schema.NewPlugin) ([]schema.Plugin, error) {
	l.Lock()
	defer l.Unlock()

	var plugins []schema.Plugin

	{
		for _, pluginGenerator := range pluginGenerators {
			plugin, pluginErr := pluginGenerator(&schema.Cfg{Viper: l.viper, Logger: l.logger})
			if pluginErr != nil {
				return plugins, pluginErr
			}
			plugins = append(plugins, plugin)
		}

	}

	{
		for _, plugin := range plugins {
			initErr := plugin.Init()
			if initErr != nil {
				l.logger.Infof("%v->Init failed:%v\n", plugin.Name(), initErr)
				return nil, initErr
			}
		}
	}
	return plugins, nil
}
