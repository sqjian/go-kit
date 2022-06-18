package aliyun

type Option interface {
	apply(*OssManager)
}

type optionFunc func(*OssManager)

func (f optionFunc) apply(s3 *OssManager) {
	f(s3)
}

func WithAddr(addr string) Option {
	return optionFunc(func(cli *OssManager) {
		cli.meta.addr = addr
	})
}

func WithKey(key string) Option {
	return optionFunc(func(cli *OssManager) {
		cli.meta.key = key
	})
}

func WithSecret(secret string) Option {
	return optionFunc(func(cli *OssManager) {
		cli.meta.secret = secret
	})
}
