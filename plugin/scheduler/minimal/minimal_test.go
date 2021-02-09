package minimal_test

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/preset/go_native/enter"
	"github.com/sqjian/go-kit/plugin/preset/go_native/leave"
	"github.com/sqjian/go-kit/plugin/proto"
	"github.com/sqjian/go-kit/plugin/scheduler/minimal"
	"github.com/sqjian/go-kit/plugin/schema"
	"testing"
)

func Test_minimal(t *testing.T) {
	minimalInst, minimalInstErr := minimal.NewMinimal(&minimal.Cfg{
		Viper:              viper.New(),
		Logger:             log.DebugLogger,
		GoNativeGenerators: []schema.NewPlugin{enter.NewPlugin, leave.NewPlugin},
	})
	if minimalInstErr != nil {
		t.Fatal(minimalInstErr)
	}

	msg := &proto.Msg{
		DataList: []*proto.Data{
			{
				Id:   "uuid",
				Desc: map[string][]byte{"key": []byte("val")},
				Data: []byte("data"),
			},
		},
	}

	_, processErr := minimalInst.Process(minimal.Dag{Id: "xxx", Steps: []string{"enter"}}, msg)
	if processErr != nil {
		t.Fatal(processErr)
	}
}
