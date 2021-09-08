package log

import (
	"github.com/sqjian/go-kit/log/provider/dummy"
	"github.com/sqjian/go-kit/log/vars"
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
