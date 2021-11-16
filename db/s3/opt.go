package s3

type Option interface {
	apply(*S3)
}

type optionFunc func(*S3)

func (f optionFunc) apply(s3 *S3) {
	f(s3)
}

func WithAddr(addr string) Option {
	return optionFunc(func(cli *S3) {
		cli.meta.addr = addr
	})
}

func WithKey(key string) Option {
	return optionFunc(func(cli *S3) {
		cli.meta.key = key
	})
}

func WithSecret(secret string) Option {
	return optionFunc(func(cli *S3) {
		cli.meta.secret = secret
	})
}

func WithConcurrency(concurrency int) Option {
	return optionFunc(func(cli *S3) {
		cli.meta.concurrency = concurrency
	})
}

func WithDebugInfo(debug bool) Option {
	return optionFunc(func(cli *S3) {
		cli.meta.debug = debug
	})
}
