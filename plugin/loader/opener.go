package loader

import (
	"fmt"
	"github.com/sqjian/go-kit/plugin/schema"
	"github.com/sqjian/go-kit/plugin/vars"
	"plugin"
	"strings"
)

func openGoPlugin(pluginNames ...string) ([]schema.NewPlugin, error) {
	reWritePluginName := func(pluginName string) string {
		if strings.HasSuffix(pluginName, ".so") {
			return pluginName
		}
		return pluginName + ".so"
	}
	verify := func(pluginName string) (schema.NewPlugin, error) {
		pg, err := plugin.Open(pluginName)
		if err != nil {
			return nil, fmt.Errorf("can not open plugin: %v,err:%w", pluginName, err)
		}

		f, err := pg.Lookup("NewPlugin")
		if err != nil {
			return nil, fmt.Errorf("lookup plugin: %v failed,err:%w", pluginName, err)
		}
		NewPlugin, NewPluginOk := f.(schema.NewPlugin)
		if !NewPluginOk {
			return nil, fmt.Errorf("%v->convert %T to %T failed,err:%w", pluginName, f, schema.NewPluginObj, vars.ErrWrapper(vars.PluginMethodNotFound))
		}
		return NewPlugin, nil
	}
	load := func(pluginNames ...string) ([]schema.NewPlugin, error) {
		if len(pluginNames) == 0 || func() bool {
			for _, pluginName := range pluginNames {
				if len(pluginName) == 0 {
					return true
				}
			}
			return false
		}() {
			return nil, fmt.Errorf("can't load plugin, please specify pluginNames")
		}

		var NewPlugins []schema.NewPlugin
		for _, pluginName := range pluginNames {
			NewPlugin, NewPluginErr := verify(reWritePluginName(pluginName))
			if NewPluginErr != nil {
				return nil, NewPluginErr
			}
			NewPlugins = append(NewPlugins, NewPlugin)
		}
		return NewPlugins, nil
	}
	return load(pluginNames...)
}
