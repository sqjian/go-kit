package container

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/schema"
	"sync"
)

type Cfg struct {
	Viper  *viper.Viper
	Logger log.API
}

type Container struct {
	sync.RWMutex
	viper   *viper.Viper
	logger  log.API
	plugins []schema.Plugin
}

func NewContainer(cfg *Cfg) Container {
	return Container{viper: cfg.Viper, logger: cfg.Logger}
}
func (p *Container) Init() error {
	return nil
}
func (p *Container) FInit() error {
	return nil
}
func (p *Container) Add(plugins []schema.Plugin) error {
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
func (p *Container) Get(pluginNames ...string) ([]schema.Plugin, error) {
	if len(pluginNames) == 0 {
		return nil, fmt.Errorf("please specify pluginNames")
	}

	p.RLock()
	defer p.RUnlock()

	pluginMap := func() map[string]schema.Plugin {
		pluginMap := make(map[string]schema.Plugin)
		for _, plugin := range p.plugins {
			pluginMap[plugin.Name()] = plugin
		}
		return pluginMap
	}()

	plugins, pluginsErr := func() ([]schema.Plugin, error) {
		var plugins []schema.Plugin
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
