package container_test

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/container"
	"github.com/sqjian/go-kit/plugin/loader"
	"github.com/sqjian/go-kit/plugin/plugin/go_native/enter"
	"github.com/sqjian/go-kit/plugin/plugin/go_native/leave"
	"testing"
)

func Test_Container(t *testing.T) {

	c := func() *container.Container {
		c := container.NewContainer(func(cfg *container.Config) {
			cfg.Viper = viper.New()
			cfg.Logger = log.TerminalLogger{}
		})
		if err := c.Init(); err != nil {
			t.Fatal(err)
		}
		return c
	}()

	l := func() *loader.Loader {
		l := loader.NewLoader(func(cfg *loader.Cfg) {
			cfg.Viper = viper.New()
			cfg.Logger = log.TerminalLogger{}
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
