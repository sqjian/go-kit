package enter

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plug/proto"
	"github.com/sqjian/go-kit/plug/schema"
)

type enter struct {
	logger log.API
	viper  *viper.Viper
}

func (e *enter) Init() error {
	e.logger.Infof("%v Init...\n", e.Name())
	return nil
}
func (e *enter) FInit() error {
	e.logger.Infof("%v FInit...\n", e.Name())
	return nil
}
func (e *enter) Interest(*proto.Msg, schema.PlugTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func (e *enter) PreProcess(msg *proto.Msg, plugTools schema.PlugTools) (schema.Action, error) {
	dagId, _ := plugTools.Get("dag.id")
	e.logger.Infof("id:%v,enter PreProcess msg:%v\n", dagId, msg.String())
	return schema.CONTINUE, nil
}

func (e *enter) Process(msg *proto.Msg, plugTools schema.PlugTools) (schema.Action, error) {
	dagId, _ := plugTools.Get("dag.id")
	e.logger.Infof("id:%v,enter Process msg:%v\n", dagId, msg.String())
	return schema.CONTINUE, nil
}

func (e *enter) PostProcess(msg *proto.Msg, plugTools schema.PlugTools) (schema.Action, error) {
	dagId, _ := plugTools.Get("dag.id")
	e.logger.Infof("id:%v,enter PostProcess msg:%v\n", dagId, msg.String())
	return schema.CONTINUE, nil
}

func (e *enter) Name() string {
	return "enter"
}

func NewPlug(fn func(cfg *schema.Cfg)) (schema.Plug, error) {
	cfg := &schema.Cfg{}
	fn(cfg)

	return &enter{viper: cfg.Viper, logger: cfg.Logger}, nil
}
