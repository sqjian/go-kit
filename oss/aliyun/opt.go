package aliyun

type Option interface {
	apply(*S3Manager)
}

type optionFunc func(*S3Manager)

func (f optionFunc) apply(s3 *S3Manager) {
	f(s3)
}

func WithAddr(addr string) Option {
	return optionFunc(func(cli *S3Manager) {
		cli.meta.addr = addr
	})
}

func WithKey(key string) Option {
	return optionFunc(func(cli *S3Manager) {
		cli.meta.key = key
	})
}

func WithSecret(secret string) Option {
	return optionFunc(func(cli *S3Manager) {
		cli.meta.secret = secret
	})
}
