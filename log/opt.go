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
		log.metaData.FileName = FileName
	})
}

func WithMaxSize(MaxSize int) Option {
	return optionFunc(func(log *logger) {
		log.metaData.MaxSize = MaxSize
	})
}

func WithMaxBackups(MaxBackups int) Option {
	return optionFunc(func(log *logger) {
		log.metaData.MaxBackups = MaxBackups
	})
}

func WithMaxAge(MaxAge int) Option {
	return optionFunc(func(log *logger) {
		log.metaData.MaxAge = MaxAge
	})
}

func WithLevel(Level Level) Option {
	return optionFunc(func(log *logger) {
		log.metaData.Level = Level
	})
}

func WithConsole(Console bool) Option {
	return optionFunc(func(log *logger) {
		log.metaData.Console = Console
	})
}

func WithLogType(logType KeyType) Option {
	return optionFunc(func(log *logger) {
		log.logType = logType
	})
}
