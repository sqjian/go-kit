package option

type options struct {
	key string
}

type option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(options *options) {
	f(options)
}

func withKey(val string) option {
	return optionFunc(func(options *options) {
		options.key = val
	})
}
