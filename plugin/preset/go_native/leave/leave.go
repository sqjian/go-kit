package leave

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plugin/proto"
	"github.com/sqjian/go-kit/plugin/schema"
)

type leave struct {
	logger log.API
	viper  *viper.Viper
}

func (l *leave) Init() error {
	l.logger.Infof("%v Init...\n", l.Name())
	return nil
}
func (l *leave) FInit() error {
	l.logger.Infof("%v FInit...\n", l.Name())
	return nil
}
func (l *leave) Interest(msg *proto.Msg, tools schema.PluginTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func (l *leave) PreProcess(msg *proto.Msg, tools schema.PluginTools) (schema.Action, error) {
	l.logger.Infof("leave PreProcess msg:%v\n", msg.String())
	return schema.CONTINUE, nil
}

func (l *leave) Process(msg *proto.Msg, tools schema.PluginTools) (schema.Action, error) {
	l.logger.Infof("leave Process msg:%v\n", msg.String())
	return schema.CONTINUE, nil
}

func (l *leave) PostProcess(msg *proto.Msg, tools schema.PluginTools) (schema.Action, error) {
	l.logger.Infof("leave PostProcess msg:%v\n", msg.String())
	return schema.CONTINUE, nil
}

func (l *leave) Name() string {
	return "leave"
}
func NewPlugin(cfg *schema.Cfg) (schema.Plugin, error) {
	return &leave{viper: cfg.Viper, logger: cfg.Logger}, nil
}
