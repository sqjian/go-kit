package loader

import (
	"fmt"
	"github.com/sqjian/go-kit/plugin/schema"
	"github.com/sqjian/go-kit/plugin/vars"
	"plugin"
	"strings"
)

func openGoPlugin(pluginNames ...string) ([]schema.NewPlug, error) {
	reWritePluginName := func(pluginName string) string {
		if strings.HasSuffix(pluginName, ".so") {
			return pluginName
		}
		return pluginName + ".so"
	}
	verify := func(pluginName string) (schema.NewPlug, error) {
		pg, err := plugin.Open(pluginName)
		if err != nil {
			return nil, fmt.Errorf("can not open plug: %v,err:%w", pluginName, err)
		}

		f, err := pg.Lookup("NewPlug")
		if err != nil {
			return nil, fmt.Errorf("lookup plug: %v failed,err:%w", pluginName, err)
		}
		NewPlugin, NewPluginOk := f.(schema.NewPlug)
		if !NewPluginOk {
			return nil, fmt.Errorf("%v->convert %T to %T failed,err:%w", pluginName, f, schema.NewPlugObj, vars.ErrWrapper(vars.PluginMethodNotFound))
		}
		return NewPlugin, nil
	}
	load := func(pluginNames ...string) ([]schema.NewPlug, error) {
		if len(pluginNames) == 0 || func() bool {
			for _, pluginName := range pluginNames {
				if len(pluginName) == 0 {
					return true
				}
			}
			return false
		}() {
			return nil, fmt.Errorf("can't load plug, please specify pluginNames")
		}

		var NewPlugins []schema.NewPlug
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
