package aliyun

type OptionFunc func(*OssManager)

func WithAddr(addr string) OptionFunc {
	return func(o *OssManager) {
		o.config.addr = addr
	}
}

func WithKey(key string) OptionFunc {
	return func(o *OssManager) {
		o.config.key = key
	}
}

func WithSecret(secret string) OptionFunc {
	return func(o *OssManager) {
		o.config.secret = secret
	}
}
