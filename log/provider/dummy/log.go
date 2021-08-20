package dummy

import (
	"encoding/json"
	"github.com/sqjian/go-kit/log/preset"
	"go.uber.org/zap"
)

func NewLogger(opts ...Option) (*logger, error) {

	loggerInst := new(logger)

	for _, opt := range opts {
		opt.apply(loggerInst)
	}

	return loggerInst, nil
}

type logger struct {
	metaData struct {
		FileName   string       /*日志的名字*/
		MaxSize    int          /*日志大小，单位MB*/
		MaxBackups int          /*日志备份个数*/
		MaxAge     int          /*日志备份时间，单位Day*/
		Level      preset.Level /*日志级别，可选：none、debug、info、warn、error*/
		Console    bool         /*是否向控制台输出*/
	}

	ready bool

	Logger *zap.SugaredLogger
}

func (l *logger) String() string {

	m := make(map[string]interface{})
	m["metaData"] = l.metaData
	res, _ := json.Marshal(m)

	return string(res)
}

func (l *logger) Debugf(_ string, _ ...interface{}) {}
func (l *logger) Debugw(_ string, _ ...interface{}) {}

func (l *logger) Infof(_ string, _ ...interface{}) {}
func (l *logger) Infow(_ string, _ ...interface{}) {}

func (l *logger) Warnf(_ string, _ ...interface{})  {}
func (l *logger) Warnw(_ string, _ ...interface{})  {}
func (l *logger) Errorf(_ string, _ ...interface{}) {}
func (l *logger) Errorw(_ string, _ ...interface{}) {}
