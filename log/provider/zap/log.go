package zap

import (
	"encoding/json"
	"fmt"
	"github.com/sqjian/go-kit/log/vars"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync/atomic"
)

func NewLogger(opts ...Option) (*logger, error) {

	loggerInst := new(logger)

	for _, opt := range opts {
		opt.apply(loggerInst)
	}

	switch {
	case len(loggerInst.metaData.FileName) == 0:
		{
			return nil, fmt.Errorf("empty fileName")
		}
	case loggerInst.metaData.MaxSize == 0:
		{
			return nil, fmt.Errorf("empty MaxSize")
		}
	case loggerInst.metaData.MaxBackups == 0:
		{
			return nil, fmt.Errorf("empty MaxBackups")
		}
	case loggerInst.metaData.MaxAge == 0:
		{
			return nil, fmt.Errorf("empty MaxAge")
		}
	case loggerInst.metaData.Level == vars.UnknownLevel:
		{
			return nil, fmt.Errorf("empty Level")
		}
	}

	err := loggerInst.init()
	if err != nil {
		return nil, err
	}
	loggerInst.Errorf("init params:%v", loggerInst)

	return loggerInst, nil
}

type logger struct {
	metaData struct {
		FileName   string     /*日志的名字*/
		MaxSize    int        /*日志大小，单位MB*/
		MaxBackups int        /*日志备份个数*/
		MaxAge     int        /*日志备份时间，单位Day*/
		Level      vars.Level /*日志级别，可选：none、debug、info、warn、error*/
		Console    bool       /*是否向控制台输出*/
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

func (l *logger) SetLevelOTF(Level vars.Level) error {

	if !l.ready {
		return fmt.Errorf("logger not ready,please init first")
	}

	atomic.StoreInt64((*int64)(&l.metaData.Level), int64(Level))

	l.Errorf("reset the level,params:%v", l)

	return nil
}

func (l *logger) init() (err error) {

	defer func() {
		if err == nil {
			l.ready = true
		}
	}()

	userPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if l.metaData.Level == vars.None {
			return false
		}
		return lvl >= func() zapcore.Level {
			switch l.metaData.Level {
			case vars.Debug:
				{
					return zapcore.DebugLevel
				}
			case vars.Info:
				{
					return zapcore.InfoLevel
				}
			case vars.Warn:
				{
					return zapcore.WarnLevel
				}
			case vars.Error:
				{
					return zapcore.ErrorLevel
				}
			default:
				{
					return zapcore.ErrorLevel
				}
			}
		}()
	})

	fileLogRotateUserWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   l.metaData.FileName,
		MaxSize:    l.metaData.MaxSize,
		MaxBackups: l.metaData.MaxBackups,
		MaxAge:     l.metaData.MaxAge,
	})

	consoleWriter := zapcore.Lock(os.Stdout)

	presetEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		NameKey:        "logger",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	var core zapcore.Core
	switch {
	case l.metaData.Console:
		{
			core = zapcore.NewTee(
				zapcore.NewCore(presetEncoder, fileLogRotateUserWriter, userPriority),
				zapcore.NewCore(presetEncoder, consoleWriter, userPriority),
			)
		}
	default:
		{
			core = zapcore.NewTee(
				zapcore.NewCore(presetEncoder, fileLogRotateUserWriter, userPriority),
			)
		}
	}

	l.Logger = zap.New(
		core,
		/*zap.AddCaller(),*/
		zap.AddCallerSkip(1),
		zap.Fields(
			zapcore.Field{
				Key:     "pid",
				Type:    zapcore.Int64Type,
				Integer: int64(os.Getpagesize()),
			},
		),
	).Sugar()

	return nil
}

func (l *logger) Debugf(template string, args ...interface{}) {
	l.Logger.Debugf(template, args...)
}
func (l *logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.Logger.Debugw(msg, keysAndValues...)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.Logger.Infof(template, args...)
}
func (l *logger) Infow(msg string, keysAndValues ...interface{}) {
	l.Logger.Infow(msg, keysAndValues...)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.Logger.Warnf(template, args...)
}
func (l *logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.Logger.Warnw(msg, keysAndValues...)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.Logger.Errorf(template, args...)
}
func (l *logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.Logger.Errorw(msg, keysAndValues...)
}
