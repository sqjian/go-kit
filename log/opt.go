package log

import (
	"fmt"
	"os"
	"strings"
)

type OptionFunc func(*config)

func WithFileName(fileName string) OptionFunc {
	return func(m *config) {
		m.FileName = fileName
	}
}

func WithCaller(caller bool, callerSkip int) OptionFunc {
	return func(m *config) {
		m.Caller = caller
		m.CallerSkip = callerSkip
	}
}

func WithMaxSize(maxSize int) OptionFunc {
	return func(m *config) {
		m.MaxSize = maxSize
	}
}

func WithMaxBackups(maxBackups int) OptionFunc {
	return func(m *config) {
		m.MaxBackups = maxBackups
	}
}

func WithMaxAge(maxAge int) OptionFunc {
	return func(m *config) {
		m.MaxAge = maxAge
	}
}

func WithLevel(level string) OptionFunc {
	return func(m *config) {
		switch strings.ToLower(level) {
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
				fmt.Printf("illegal log level:%v", level)
				os.Exit(-1)
			}
		}
	}
}

func WithConsole(console bool) OptionFunc {
	return func(m *config) {
		m.Console = console
	}
}
