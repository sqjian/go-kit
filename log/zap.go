package log

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

func GenZapLog(
	FileName string,
	MaxSize int /*megabytes*/,
	MaxBackups int,
	MaxAge int /*days*/,
	Level string,
	Console bool, /*enable write to console*/
) (*zap.SugaredLogger, error) {
	return genZapLog(FileName, MaxSize, MaxBackups, MaxAge, Level, Console)
}

func genZapLog(
	FileName string,
	MaxSize int /*megabytes*/,
	MaxBackups int,/*number of backup*/
	MaxAge int /*days*/,
	Level string,
	Console bool,
) (*zap.SugaredLogger, error) {
	if Level != "none" &&
		Level != "debug" &&
		Level != "info" &&
		Level != "warn" &&
		Level != "error" {
		return nil, fmt.Errorf("logging level:%v is illegal", Level)
	}

	userPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		if Level == "none" {
			return false
		}
		return lvl >= func() zapcore.Level {
			switch Level {
			case "debug":
				{
					return zapcore.DebugLevel
				}
			case "info":
				{
					return zapcore.InfoLevel
				}
			case "warn":
				{
					return zapcore.WarnLevel
				}
			case "error":
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
		Filename:   FileName,
		MaxSize:    MaxSize,
		MaxBackups: MaxBackups,
		MaxAge:     MaxAge,
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
	case Console:
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

	return zap.New(
		core,
		zap.AddCaller(),
		zap.Fields(
			zapcore.Field{
				Key:     "pid",
				Type:    zapcore.Int64Type,
				Integer: int64(os.Getpagesize()),
			},
		),
	).Sugar(), nil
}
