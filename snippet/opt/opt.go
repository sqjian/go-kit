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

type OptionFunc func(*options)

func WithKey(val string) OptionFunc {
	return func(options *options) {
		options.key = val
	}
}
