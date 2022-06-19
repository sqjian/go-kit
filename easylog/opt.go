package easylog

type Option interface {
	apply(*Meta)
}

type optionFunc func(*Meta)

func (f optionFunc) apply(log *Meta) {
	f(log)
}

func WithFileName(FileName string) Option {
	return optionFunc(func(m *Meta) {
		m.FileName = FileName
	})
}

func WithCaller(Caller bool, CallerSkip int) Option {
	return optionFunc(func(m *Meta) {
		m.Caller = Caller
		m.CallerSkip = CallerSkip
	})
}

func WithMaxSize(MaxSize int) Option {
	return optionFunc(func(m *Meta) {
		m.MaxSize = MaxSize
	})
}

func WithMaxBackups(MaxBackups int) Option {
	return optionFunc(func(m *Meta) {
		m.MaxBackups = MaxBackups
	})
}

func WithMaxAge(MaxAge int) Option {
	return optionFunc(func(m *Meta) {
		m.MaxAge = MaxAge
	})
}

func WithLevel(Level Level) Option {
	return optionFunc(func(m *Meta) {
		m.Level = Level
	})
}

func WithConsole(Console bool) Option {
	return optionFunc(func(m *Meta) {
		m.Console = Console
	})
}

func WithBuilder(builder func(*Meta) (API, error)) Option {
	return optionFunc(func(m *Meta) {
		m.builder = builder
	})
}
