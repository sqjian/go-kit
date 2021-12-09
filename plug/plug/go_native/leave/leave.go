package leave

import (
	"github.com/spf13/viper"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/plug/proto"
	"github.com/sqjian/go-kit/plug/schema"
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
func (l *leave) Interest(*proto.Msg, schema.PlugTools) (schema.Action, error) {
	return schema.CONTINUE, nil
}

func (l *leave) PreProcess(msg *proto.Msg, plugTools schema.PlugTools) (schema.Action, error) {
	dagId, _ := plugTools.Get("dag.id")
	l.logger.Infof("id:%v,leave PreProcess msg:%v\n", dagId, msg.String())
	return schema.CONTINUE, nil
}

func (l *leave) Process(msg *proto.Msg, plugTools schema.PlugTools) (schema.Action, error) {
	dagId, _ := plugTools.Get("dag.id")
	l.logger.Infof("id:%v,leave Process msg:%v\n", dagId, msg.String())
	return schema.CONTINUE, nil
}

func (l *leave) PostProcess(msg *proto.Msg, plugTools schema.PlugTools) (schema.Action, error) {
	dagId, _ := plugTools.Get("dag.id")
	l.logger.Infof("id:%v,leave PostProcess msg:%v\n", dagId, msg.String())
	return schema.CONTINUE, nil
}

func (l *leave) Name() string {
	return "leave"
}
func NewPlug(fn func(cfg *schema.Cfg)) (schema.Plug, error) {
	cfg := &schema.Cfg{}
	fn(cfg)

	return &leave{viper: cfg.Viper, logger: cfg.Logger}, nil
}
