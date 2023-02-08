package exec

import (
	"io"
)

type OptionFunc func(*Config)

func WithArgs(a ...string) OptionFunc {
	return OptionFunc(func(c *Config) {
		c.args = a
	})
}

func WithWriters(w ...io.WriteCloser) OptionFunc {
	return OptionFunc(func(c *Config) {
		c.writers = w
	})
}
