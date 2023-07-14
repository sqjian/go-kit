package log

type OptionFunc func(*config)

func WithFileName(FileName string) OptionFunc {
	return func(m *config) {
		m.FileName = FileName
	}
}

func WithCaller(Caller bool, CallerSkip int) OptionFunc {
	return func(m *config) {
		m.Caller = Caller
		m.CallerSkip = CallerSkip
	}
}

func WithMaxSize(MaxSize int) OptionFunc {
	return func(m *config) {
		m.MaxSize = MaxSize
	}
}

func WithMaxBackups(MaxBackups int) OptionFunc {
	return func(m *config) {
		m.MaxBackups = MaxBackups
	}
}

func WithMaxAge(MaxAge int) OptionFunc {
	return func(m *config) {
		m.MaxAge = MaxAge
	}
}

func WithLevel(Level Level) OptionFunc {
	return func(m *config) {
		m.Level = Level
	}
}

func WithConsole(Console bool) OptionFunc {
	return func(m *config) {
		m.Console = Console
	}
}
