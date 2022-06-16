package log

var (
	DebugLogger = func() API {
		logger, _ := NewLogger(
			WithFileName("go-kit.log"),
			WithMaxSize(3),
			WithMaxBackups(3),
			WithMaxAge(3),
			WithLevel(Dummy),
			WithConsole(true),
			WithCaller(false),
		)
		return logger
	}()
	DummyLogger = func() API {
		logger, _ := NewLogger(
			WithFileName("go-kit.log"),
			WithMaxSize(3),
			WithMaxBackups(3),
			WithMaxAge(3),
			WithLevel(Dummy),
			WithConsole(false),
		)
		return logger
	}()
)

func Debugf(template string, args ...interface{}) {
	DebugLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	DebugLogger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	DebugLogger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	DebugLogger.Errorf(template, args...)
}
