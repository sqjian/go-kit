package enter

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/proto"
	"github.com/sqjian/go-kit/plugin/schema"
)

type enter struct {
	logger log.Logger
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
func (e *enter) Interest(msg *proto.Msg, tools schema.PluginTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func (e *enter) PreProcess(msg *proto.Msg, tools schema.PluginTools) (schema.Action, error) {
	e.logger.Infof("enter PreProcess msg:%v\n", msg.String())
	return schema.CONTINUE, nil
}

func (e *enter) Process(msg *proto.Msg, tools schema.PluginTools) (schema.Action, error) {
	e.logger.Infof("enter Process msg:%v\n", msg.String())
	return schema.CONTINUE, nil
}

func (e *enter) PostProcess(msg *proto.Msg, tools schema.PluginTools) (schema.Action, error) {
	e.logger.Infof("enter PostProcess msg:%v\n", msg.String())
	return schema.CONTINUE, nil
}

func (e *enter) Name() string {
	return "enter"
}

func NewPlugin(cfg *schema.Cfg) (schema.Plugin, error) {
	return &enter{viper: cfg.Viper, logger: cfg.Logger}, nil
}
