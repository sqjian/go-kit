package log

type OptionFunc func(*Meta)

func WithFileName(FileName string) OptionFunc {
	return func(m *Meta) {
		m.FileName = FileName
	}
}

func WithCaller(Caller bool, CallerSkip int) OptionFunc {
	return func(m *Meta) {
		m.Caller = Caller
		m.CallerSkip = CallerSkip
	}
}

func WithMaxSize(MaxSize int) OptionFunc {
	return func(m *Meta) {
		m.MaxSize = MaxSize
	}
}

func WithMaxBackups(MaxBackups int) OptionFunc {
	return func(m *Meta) {
		m.MaxBackups = MaxBackups
	}
}

func WithMaxAge(MaxAge int) OptionFunc {
	return func(m *Meta) {
		m.MaxAge = MaxAge
	}
}

func WithLevel(Level Level) OptionFunc {
	return func(m *Meta) {
		m.Level = Level
	}
}

func WithConsole(Console bool) OptionFunc {
	return func(m *Meta) {
		m.Console = Console
	}
}

func WithBuilder(builder func(*Meta) (API, error)) OptionFunc {
	return func(m *Meta) {
		m.builder = builder
	}
}
