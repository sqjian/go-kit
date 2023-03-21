package log

type OptionFunc func(*Config)

func WithFileName(FileName string) OptionFunc {
	return func(m *Config) {
		m.FileName = FileName
	}
}

func WithCaller(Caller bool, CallerSkip int) OptionFunc {
	return func(m *Config) {
		m.Caller = Caller
		m.CallerSkip = CallerSkip
	}
}

func WithMaxSize(MaxSize int) OptionFunc {
	return func(m *Config) {
		m.MaxSize = MaxSize
	}
}

func WithMaxBackups(MaxBackups int) OptionFunc {
	return func(m *Config) {
		m.MaxBackups = MaxBackups
	}
}

func WithMaxAge(MaxAge int) OptionFunc {
	return func(m *Config) {
		m.MaxAge = MaxAge
	}
}

func WithLevel(Level Level) OptionFunc {
	return func(m *Config) {
		m.Level = Level
	}
}

func WithConsole(Console bool) OptionFunc {
	return func(m *Config) {
		m.Console = Console
	}
}

func WithBuilder(builder func(*Config) (API, error)) OptionFunc {
	return func(m *Config) {
		m.builder = builder
	}
}
