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
	defaultLogType = vars.Dummy
)

var (
	DummyLogger = func() Logger { logger, _ := dummy.NewLogger(); return logger }()
	DebugLogger = func() Logger {
		logger, _ := NewLogger(
			WithFileName("go-kit.log"),
			WithMaxSize(3),
			WithMaxBackups(3),
			WithMaxAge(3),
			WithLevel(vars.Debug),
			WithLogType(vars.Zap),
			WithConsole(true),
		)
		return logger
	}()
)

type logger struct {
	logType  vars.LogType
	metaData struct {
		FileName   string     /*日志的名字*/
		MaxSize    int        /*日志大小，单位MB*/
		MaxBackups int        /*日志备份个数*/
		MaxAge     int        /*日志备份时间，单位Day*/
		Level      vars.Level /*日志级别，可选：none、debug、info、warn、error*/
		Console    bool       /*是否向控制台输出*/
	}
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

	{
		switch {
		case len(loggerInst.metaData.FileName) == 0:
			fallthrough
		case loggerInst.metaData.MaxSize == 0:
			fallthrough
		case loggerInst.metaData.MaxBackups == 0:
			fallthrough
		case loggerInst.metaData.MaxAge == 0:
			fallthrough
		case loggerInst.metaData.Level == vars.UnknownLevel:
			return nil, vars.ErrWrapper(vars.IllegalParams)
		}
	}

	switch loggerInst.logType {
	case vars.Zap:
		{
			return zap.NewLogger(
				zap.WithFileName(loggerInst.metaData.FileName),
				zap.WithMaxSize(loggerInst.metaData.MaxSize),
				zap.WithMaxBackups(loggerInst.metaData.MaxBackups),
				zap.WithMaxAge(loggerInst.metaData.MaxAge),
				zap.WithLevel(loggerInst.metaData.Level),
				zap.WithConsole(loggerInst.metaData.Console),
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
