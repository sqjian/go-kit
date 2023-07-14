package container

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/schema"
	"sync"
)

type Config struct {
	Viper  *viper.Viper
	Logger log.Log
}

type Container struct {
	sync.RWMutex
	viper   *viper.Viper
	logger  log.Log
	plugins []schema.Plug
}

func NewContainer(fn func(*Config)) *Container {
	cfg := &Config{}
	fn(cfg)

	return &Container{viper: cfg.Viper, logger: cfg.Logger}
}
func (p *Container) Init() error {
	return nil
}
func (p *Container) FInit() error {
	return nil
}
func (p *Container) Add(plugins []schema.Plug) error {
	p.Lock()
	defer p.Unlock()

	checkRepeatAdd := func(pluginName string) bool {
		for _, plugin := range p.plugins {
			if plugin.Name() == pluginName {
				return true
			}
		}
		return false
	}
	for _, plugin := range plugins {
		if checkRepeatAdd(plugin.Name()) {
			return fmt.Errorf("%v already added", plugin.Name())
		}
		p.plugins = append(p.plugins, plugin)
	}

	return nil
}

func (p *Container) Remove(pluginName string) error {
	p.Lock()
	defer p.Unlock()
	for ix, val := range p.plugins {
		if val.Name() == pluginName {
			if ix+1 >= len(p.plugins) {
				p.plugins = p.plugins[:ix]
			} else {
				p.plugins = append(p.plugins[:ix], p.plugins[ix+1:]...)
			}
			return p.plugins[ix].FInit()
		}
	}
	return nil
}
func (p *Container) Get(pluginNames ...string) ([]schema.Plug, error) {
	if len(pluginNames) == 0 {
		return nil, fmt.Errorf("please specify pluginNames")
	}

	{
		checkRepeatAdd := func(pluginName string) bool {
			cnt := 0
			for _, plugin := range pluginNames {
				if plugin == pluginName {
					cnt++
				}
			}
			if cnt > 1 {
				return true
			}
			return false
		}
		for _, plugName := range pluginNames {
			if checkRepeatAdd(plugName) {
				return nil, fmt.Errorf("plug->%v repeat, all plugs:%v", plugName, pluginNames)
			}
		}
	}

	p.RLock()
	defer p.RUnlock()

	pluginMap := func() map[string]schema.Plug {
		pluginMap := make(map[string]schema.Plug)
		for _, plugin := range p.plugins {
			pluginMap[plugin.Name()] = plugin
		}
		return pluginMap
	}()

	plugins, pluginsErr := func() ([]schema.Plug, error) {
		var plugins []schema.Plug
		for _, pluginName := range pluginNames {
			plugin, pluginOk := pluginMap[pluginName]
			if !pluginOk {
				return nil, fmt.Errorf("can not found %v", pluginName)
			}
			plugins = append(plugins, plugin)
		}
		return plugins, nil
	}()

	return plugins, pluginsErr
}
