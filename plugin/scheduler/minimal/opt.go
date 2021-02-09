package minimal

func newDefaultOptions() *options {
	return &options{kvs: make(map[string]interface{})}
}

type options struct {
	kvs map[string]interface{}
}

type Option interface {
	apply(*options)
}

type optionFunc func(*options)

func (f optionFunc) apply(options *options) {
	f(options)
}

func WithKvs(kvs map[string]interface{}) Option {
	return optionFunc(func(in *options) {
		in.kvs = kvs
	})
}
