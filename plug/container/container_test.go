package container_test

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/easylog"
	"github.com/sqjian/go-kit/plug/container"
	"github.com/sqjian/go-kit/plug/loader"
	"github.com/sqjian/go-kit/plug/plug/go_native/enter"
	"github.com/sqjian/go-kit/plug/plug/go_native/leave"
	"testing"
)

func Test_Container(t *testing.T) {

	c := func() *container.Container {
		c := container.NewContainer(func(cfg *container.Cfg) {
			cfg.Viper = viper.New()
			cfg.Logger = easylog.DebugLogger
		})
		if err := c.Init(); err != nil {
			t.Fatal(err)
		}
		return c
	}()

	l := func() *loader.Loader {
		l := loader.NewLoader(func(cfg *loader.Cfg) {
			cfg.Viper = viper.New()
			cfg.Logger = easylog.DebugLogger
		})
		if err := l.Init(); err != nil {
			t.Fatal(err)
		}
		return l
	}()

	{
		plugins, pluginsErr := l.Load(enter.NewPlug, leave.NewPlug)
		if pluginsErr != nil {
			t.Fatal(pluginsErr)
		}

		addErr := c.Add(plugins)
		if addErr != nil {
			t.Fatal(addErr)
		}

		plugins, pluginsErr = c.Get("enter")
		if pluginsErr != nil {
			t.Fatal(pluginsErr)
		}
		for ix, plugin := range plugins {
			t.Logf("test->plugin_%v:%v\n", ix, plugin.Name())
		}
	}

	if err := c.FInit(); err != nil {
		t.Fatal(err)
	}
}
