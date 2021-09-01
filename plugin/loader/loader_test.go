package loader

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/preset/go_native/enter"
	"github.com/sqjian/go-kit/plugin/preset/go_native/leave"
	"testing"
)

func Test_Loader(t *testing.T) {
	loader := NewLoader(&Cfg{viper.New(), log.DebugLogger})
	if err := loader.Init(); err != nil {
		t.Fatal(err)
	}
	plugins, pluginsErr := loader.Load(enter.NewPlugin, leave.NewPlugin)
	if pluginsErr != nil {
		t.Fatal(pluginsErr)
	}
	for ix, plugin := range plugins {
		t.Logf("test->plugin_%v:%v\n", ix, plugin.Name())
	}
	if err := loader.FInit(); err != nil {
		t.Fatal(err)
	}
}