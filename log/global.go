package log

var (
	DebugLogger = func() API {
		logger, _ := NewLogger(
			WithFileName("go-kit.log"),
			WithMaxSize(3),
			WithMaxBackups(3),
			WithMaxAge(3),
			WithLevel(Debug),
			WithConsole(true),
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
			WithConsole(true),
		)
		return logger
	}()
)
