package log

import (
	"encoding/json"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync/atomic"
)

func newZapLogger(config *config) *zapLogger {
	userFilePriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if config.Level == Dummy {
			return false
		}
		return lvl >= func() zapcore.Level {
			switch config.Level {
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
		Filename:   config.FileName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
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
	case config.Console:
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

	logInst := zap.New(
		core,
		zap.WithCaller(config.Caller),
		zap.AddCallerSkip(config.CallerSkip),
		zap.Fields(
			zapcore.Field{
				Key:     "pid",
				Type:    zapcore.Int64Type,
				Integer: int64(os.Getpagesize()),
			},
		),
	).Sugar()

	return &zapLogger{
		config:        config,
		SugaredLogger: logInst,
	}

}

type zapLogger struct {
	config *config
	*zap.SugaredLogger
}

func (l *zapLogger) String() string {
	res, _ := json.Marshal(l.config)
	return string(res)
}

func (l *zapLogger) SetLevelOTF(Level Level) error {
	atomic.StoreInt64((*int64)(&l.config.Level), int64(Level))

	l.Errorf("reset the level,params:%v", l)

	return nil
}

func (l *zapLogger) Debugf(template string, args ...any) {
	l.SugaredLogger.Debugf(template, args...)
}

func (l *zapLogger) Infof(template string, args ...any) {
	l.SugaredLogger.Infof(template, args...)
}

func (l *zapLogger) Warnf(template string, args ...any) {
	l.SugaredLogger.Warnf(template, args...)
}

func (l *zapLogger) Errorf(template string, args ...any) {
	l.SugaredLogger.Errorf(template, args...)
}
