package container_test

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/container"
	"github.com/sqjian/go-kit/plugin/loader"
	"github.com/sqjian/go-kit/plugin/preset/go_native/enter"
	"github.com/sqjian/go-kit/plugin/preset/go_native/leave"
	"testing"
)

func Test_Container(t *testing.T) {
	container := container.NewContainer(&container.Cfg{Viper: viper.New(), Logger: log.DebugLogger})
	if err := container.Init(); err != nil {
		t.Fatal(err)
	}

	loader := loader.NewLoader(&loader.Cfg{Viper: viper.New(), Logger: log.DebugLogger})
	if err := loader.Init(); err != nil {
		t.Fatal(err)
	}
	plugins, pluginsErr := loader.Load(enter.NewPlugin, leave.NewPlugin)
	if pluginsErr != nil {
		t.Fatal(pluginsErr)
	}

	addErr := container.Add(plugins)
	if addErr != nil {
		t.Fatal(addErr)
	}

	plugins, pluginsErr = container.Get("enter")
	if pluginsErr != nil {
		t.Fatal(pluginsErr)
	}
	for ix, plugin := range plugins {
		t.Logf("test->plugin_%v:%v\n", ix, plugin.Name())
	}

	if err := container.FInit(); err != nil {
		t.Fatal(err)
	}
}
