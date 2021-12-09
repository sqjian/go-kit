package mini_test

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plug/plug/go_native/enter"
	"github.com/sqjian/go-kit/plug/plug/go_native/leave"
	"github.com/sqjian/go-kit/plug/proto"
	"github.com/sqjian/go-kit/plug/scheduler/mini"
	"github.com/sqjian/go-kit/plug/schema"
	"testing"
)

func Test_minimal(t *testing.T) {
	minimalInst, minimalInstErr := mini.NewMinimal(func(cfg *mini.Cfg) {
		cfg.Viper = viper.New()
		cfg.Logger = log.DebugLogger
		cfg.Plugs = []schema.NewPlug{enter.NewPlugin, leave.NewPlugin}
	})
	if minimalInstErr != nil {
		t.Fatal(minimalInstErr)
	}

	msg := &proto.Msg{
		Id:   "uuid",
		Desc: map[string][]byte{"key": []byte("val")},
	}
	_, processErr := minimalInst.Process(mini.Dag{Id: "xxx", Steps: []string{"enter"}}, msg)
	if processErr != nil {
		t.Fatal(processErr)
	}
}
