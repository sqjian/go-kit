package log

type Option interface {
	apply(*logger)
}

type optionFunc func(*logger)

func (f optionFunc) apply(log *logger) {
	f(log)
}

func WithFileName(FileName string) Option {
	return optionFunc(func(log *logger) {
		log.MetaData.FileName = FileName
	})
}

func WithMaxSize(MaxSize int) Option {
	return optionFunc(func(log *logger) {
		log.MetaData.MaxSize = MaxSize
	})
}

func WithMaxBackups(MaxBackups int) Option {
	return optionFunc(func(log *logger) {
		log.MetaData.MaxBackups = MaxBackups
	})
}

func WithMaxAge(MaxAge int) Option {
	return optionFunc(func(log *logger) {
		log.MetaData.MaxAge = MaxAge
	})
}

func WithLevel(Level Level) Option {
	return optionFunc(func(log *logger) {
		log.MetaData.Level = Level
	})
}

func WithConsole(Console bool) Option {
	return optionFunc(func(log *logger) {
		log.MetaData.Console = Console
	})
}
