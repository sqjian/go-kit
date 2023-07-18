package log

import "strings"

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

func WithLevel(Level string) OptionFunc {
	return func(m *config) {
		switch strings.ToLower(Level) {
		case "error":
			{
				m.Level = Error
			}
		case "warn":
			{
				m.Level = Warn
			}
		case "info":
			{
				m.Level = Info
			}
		case "debug":
			{
				m.Level = Debug
			}
		case "dummy":
			{
				m.Level = Dummy
			}
		default:
			{
				m.Level = Error
			}
		}
	}
}

func WithConsole(Console bool) OptionFunc {
	return func(m *config) {
		m.Console = Console
	}
}
