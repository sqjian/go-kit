package log

import "github.com/sqjian/go-kit/log/vars"

type Option interface {
	apply(*logger)
}

type optionFunc func(*logger)

func (f optionFunc) apply(log *logger) {
	f(log)
}

func WithFileName(FileName string) Option {
	return optionFunc(func(log *logger) {
		log.meta.FileName = FileName
	})
}

func WithMaxSize(MaxSize int) Option {
	return optionFunc(func(log *logger) {
		log.meta.MaxSize = MaxSize
	})
}

func WithMaxBackups(MaxBackups int) Option {
	return optionFunc(func(log *logger) {
		log.meta.MaxBackups = MaxBackups
	})
}

func WithMaxAge(MaxAge int) Option {
	return optionFunc(func(log *logger) {
		log.meta.MaxAge = MaxAge
	})
}

func WithLevel(Level vars.Level) Option {
	return optionFunc(func(log *logger) {
		log.meta.Level = Level
	})
}

func WithConsole(Console bool) Option {
	return optionFunc(func(log *logger) {
		log.meta.Console = Console
	})
}

func WithLogType(logType vars.LogType) Option {
	return optionFunc(func(log *logger) {
		log.logType = logType
	})
}
