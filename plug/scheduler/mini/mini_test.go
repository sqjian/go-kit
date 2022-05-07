package mini_test

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plug/proto"
	"github.com/sqjian/go-kit/plug/scheduler/mini"
	"github.com/sqjian/go-kit/plug/schema"
	"testing"
)

type Test1 struct {
	logger log.API
	viper  *viper.Viper
}

func (p Test1) Name() string {
	return "Test1"
}

func (p Test1) Init() error {
	return nil
}

func (p Test1) FInit() error {
	return nil
}

func (p Test1) Interest(msg *proto.Msg, tools schema.PlugTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func (p Test1) PreProcess(msg *proto.Msg, tools schema.PlugTools) (schema.Action, error) {
	dagId, _ := tools.Get("dag.id")
	p.logger.Infof("id:%v,enter PreProcess msg:%v\n", dagId, msg.String())
	return schema.CONTINUE, nil
}

func (p Test1) Process(msg *proto.Msg, tools schema.PlugTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func (p Test1) PostProcess(msg *proto.Msg, tools schema.PlugTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func NewTest1Plug(fn func(cfg *schema.Cfg)) (schema.Plug, error) {
	cfg := &schema.Cfg{}
	fn(cfg)

	return &Test1{viper: cfg.Viper, logger: cfg.Logger}, nil
}

type Test2 struct {
	logger log.API
	viper  *viper.Viper
}

func (p Test2) Name() string {
	return "Test2"
}

func (p Test2) Init() error {
	return nil
}

func (p Test2) FInit() error {
	return nil
}

func (p Test2) Interest(msg *proto.Msg, tools schema.PlugTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func (p Test2) PreProcess(msg *proto.Msg, tools schema.PlugTools) (schema.Action, error) {
	dagId, _ := tools.Get("dag.id")
	p.logger.Infof("id:%v,enter PreProcess msg:%v\n", dagId, msg.String())
	return schema.CONTINUE, nil
}

func (p Test2) Process(msg *proto.Msg, tools schema.PlugTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func (p Test2) PostProcess(msg *proto.Msg, tools schema.PlugTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func NewTest2Plug(fn func(cfg *schema.Cfg)) (schema.Plug, error) {
	cfg := &schema.Cfg{}
	fn(cfg)

	return &Test2{viper: cfg.Viper, logger: cfg.Logger}, nil
}

func Test_minimal(t *testing.T) {
	minimalInst, minimalInstErr := mini.NewMinimal(func(cfg *mini.Cfg) {
		cfg.Viper = viper.New()
		cfg.Logger = log.DebugLogger
		cfg.Plugs = []schema.NewPlug{NewTest1Plug, NewTest2Plug}
	})
	if minimalInstErr != nil {
		t.Fatal(minimalInstErr)
	}

	msg := &proto.Msg{
		Id:   "uuid",
		Desc: map[string][]byte{"key": []byte("val")},
	}
	_, processErr := minimalInst.Process(mini.Dag{Id: "xxx", Steps: []string{"Test2", "Test1"}}, msg)
	if processErr != nil {
		t.Fatal(processErr)
	}
}
