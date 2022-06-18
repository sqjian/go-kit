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

func newZapLogger(meta *Meta) (*zapLogger, error) {

	zapInst := &zapLogger{
		meta: meta,
	}

	err := zapInst.init()
	if err != nil {
		return nil, err
	}

	return zapInst, nil
}

type zapLogger struct {
	meta   *Meta
	ready  bool
	Logger *zap.SugaredLogger
}

func (l *zapLogger) String() string {

	m := make(map[string]interface{})
	m["meta"] = l.meta
	res, _ := json.Marshal(m)

	return string(res)
}

func (l *zapLogger) SetLevelOTF(Level Level) error {

	if !l.ready {
		return fmt.Errorf("zapLogger not ready,please init first")
	}

	atomic.StoreInt64((*int64)(&l.meta.Level), int64(Level))

	l.Errorf("reset the level,params:%v", l)

	return nil
}

func (l *zapLogger) init() (err error) {

	defer func() {
		if err == nil {
			l.ready = true
		}
	}()

	userFilePriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if l.meta.Level == Dummy {
			return false
		}
		return lvl >= func() zapcore.Level {
			switch l.meta.Level {
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
	userConsolePriority := zap.LevelEnablerFunc(func(_ zapcore.Level) bool {
		return true
	})

	fileLogRotateUserWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   l.meta.FileName,
		MaxSize:    l.meta.MaxSize,
		MaxBackups: l.meta.MaxBackups,
		MaxAge:     l.meta.MaxAge,
	})

	consoleWriter := zapcore.Lock(os.Stdout)

	presetEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		CallerKey:      "caller",
		NameKey:        "zapLogger",
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
	case l.meta.Console:
		{
			core = zapcore.NewTee(
				zapcore.NewCore(presetEncoder, fileLogRotateUserWriter, userFilePriority),
				zapcore.NewCore(presetEncoder, consoleWriter, userConsolePriority),
			)
		}
	default:
		{
			core = zapcore.NewTee(
				zapcore.NewCore(presetEncoder, fileLogRotateUserWriter, userFilePriority),
			)
		}
	}

	l.Logger = zap.New(
		core,
		zap.WithCaller(l.meta.Caller),
		zap.AddCallerSkip(l.meta.CallerSkip),
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

func (l *zapLogger) Debugf(template string, args ...interface{}) {
	l.Logger.Debugf(template, args...)
}

func (l *zapLogger) Infof(template string, args ...interface{}) {
	l.Logger.Infof(template, args...)
}

func (l *zapLogger) Warnf(template string, args ...interface{}) {
	l.Logger.Warnf(template, args...)
}

func (l *zapLogger) Errorf(template string, args ...interface{}) {
	l.Logger.Errorf(template, args...)
}
