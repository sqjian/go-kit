package log

import (
	"github.com/sqjian/go-kit/log/provider/dummy"
	"github.com/sqjian/go-kit/log/provider/zap"
	"github.com/sqjian/go-kit/log/vars"
)

type Logger interface {
	Debugf(template string, args ...interface{})
	Debugw(msg string, keysAndValues ...interface{})

	Infof(template string, args ...interface{})
	Infow(msg string, keysAndValues ...interface{})

	Warnf(template string, args ...interface{})
	Warnw(msg string, keysAndValues ...interface{})

	Errorf(template string, args ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
}

const (
	defaultLogType = vars.Zap
)

type logger struct {
	meta struct {
		FileName   string     /*日志的名字*/
		MaxSize    int        /*日志大小，单位MB*/
		MaxBackups int        /*日志备份个数*/
		MaxAge     int        /*日志备份时间，单位Day*/
		Level      vars.Level /*日志级别，可选：none、debug、info、warn、error*/
		Console    bool       /*是否向控制台输出*/
	}
	logType vars.LogType
}

func newDefaultLoggerConfig() *logger {
	return &logger{
		logType: defaultLogType,
	}
}

func NewLogger(opts ...Option) (Logger, error) {

	loggerInst := newDefaultLoggerConfig()

	for _, opt := range opts {
		opt.apply(loggerInst)
	}

	switch {
	case len(loggerInst.meta.FileName) == 0:
		fallthrough
	case loggerInst.meta.MaxSize == 0:
		fallthrough
	case loggerInst.meta.MaxBackups == 0:
		fallthrough
	case loggerInst.meta.MaxAge == 0:
		fallthrough
	case loggerInst.meta.Level == vars.UnknownLevel:
		return nil, vars.ErrWrapper(vars.IllegalParams)
	}

	switch loggerInst.logType {
	case vars.Zap:
		{
			return zap.NewLogger(
				zap.WithFileName(loggerInst.meta.FileName),
				zap.WithMaxSize(loggerInst.meta.MaxSize),
				zap.WithMaxBackups(loggerInst.meta.MaxBackups),
				zap.WithMaxAge(loggerInst.meta.MaxAge),
				zap.WithLevel(loggerInst.meta.Level),
				zap.WithConsole(loggerInst.meta.Console),
			)
		}
	case vars.Dummy:
		{
			return dummy.NewLogger()
		}
	default:
		{
			return nil, vars.ErrWrapper(vars.IllegalKeyType)
		}
	}
}
