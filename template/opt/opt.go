package opt

const (
	defaultKey = ""
)

func newDefaultOptions() *options {
	return &options{
		key: defaultKey,
	}
}

type options struct {
	key string
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(options *options) {
	f(options)
}

func withKey(val string) Option {
	return optionFunc(func(options *options) {
		options.key = val
	})
}
