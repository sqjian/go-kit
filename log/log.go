package log

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync/atomic"
)

type Level = int64

const (
	Unknown Level = iota
	None
	Debug
	Info
	Warn
	Error
)

func NewLogger(opts ...Option) (*logger, error) {

	loggerTmp := new(logger)

	for _, opt := range opts {
		opt.apply(loggerTmp)
	}

	if len(loggerTmp.MetaData.FileName) == 0 {
		return nil, fmt.Errorf("empty fileName")
	}
	if loggerTmp.MetaData.MaxSize == 0 {
		return nil, fmt.Errorf("empty MaxSize")
	}
	if loggerTmp.MetaData.MaxBackups == 0 {
		return nil, fmt.Errorf("empty MaxBackups")
	}
	if loggerTmp.MetaData.MaxAge == 0 {
		return nil, fmt.Errorf("empty MaxAge")
	}
	if loggerTmp.MetaData.Level == Unknown {
		return nil, fmt.Errorf("empty Level")
	}

	err := loggerTmp.init()
	if err != nil {
		return nil, err
	}
	loggerTmp.Errorf("init params:%v", loggerTmp)

	return loggerTmp, nil
}

type logger struct {
	MetaData struct {
		FileName   string /*日志的名字*/
		MaxSize    int    /*日志大小，单位MB*/
		MaxBackups int    /*日志备份个数*/
		MaxAge     int    /*日志备份时间，单位Day*/
		Level      Level  /*日志级别，可选：none、debug、info、warn、error*/
		Console    bool   /*是否向控制台输出*/
	}

	ready bool

	Logger *zap.SugaredLogger
}

func (l *logger) String() string {

	m := make(map[string]interface{})
	m["metaData"] = l.MetaData
	res, _ := json.Marshal(m)

	return string(res)
}

func (l *logger) SetLevelOTF(Level Level) error {

	if !l.ready {
		return fmt.Errorf("logger not ready,please init first")
	}

	atomic.StoreInt64(&l.MetaData.Level, Level)

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
		if l.MetaData.Level == None {
			return false
		}
		return lvl >= func() zapcore.Level {
			switch l.MetaData.Level {
			case Debug:
				{
					return zapcore.DebugLevel
				}
			case Info:
				{
					return zapcore.InfoLevel
				}
			case Warn:
				{
					return zapcore.WarnLevel
				}
			case Error:
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
		Filename:   l.MetaData.FileName,
		MaxSize:    l.MetaData.MaxSize,
		MaxBackups: l.MetaData.MaxBackups,
		MaxAge:     l.MetaData.MaxAge,
	})

	consoleWriter := zapcore.Lock(os.Stdout)

	commonEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
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
	case l.MetaData.Console:
		{
			core = zapcore.NewTee(
				zapcore.NewCore(commonEncoder, fileLogRotateUserWriter, userPriority),
				zapcore.NewCore(commonEncoder, consoleWriter, userPriority),
			)
		}
	default:
		{
			core = zapcore.NewTee(
				zapcore.NewCore(commonEncoder, fileLogRotateUserWriter, userPriority),
			)
		}
	}

	l.Logger = zap.New(
		core,
		zap.AddCaller(),
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
