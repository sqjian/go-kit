package mini

import "sync"

func newDefOptions() *opts {
	return &opts{kvs: sync.Map{}}
}

type opts struct {
	kvs sync.Map
}

type Opt interface {
	apply(*opts)
}

type optFn func(*opts)

func (f optFn) apply(opts *opts) {
	f(opts)
}

func WithKv(k any, v any) Opt {
	return optFn(func(in *opts) {
		in.kvs.Store(k, v)
	})
}
