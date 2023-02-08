package aliyun

type OptionFunc func(*OssManager)

func WithAddr(addr string) OptionFunc {
	return func(o *OssManager) {
		o.meta.addr = addr
	}
}

func WithKey(key string) OptionFunc {
	return func(o *OssManager) {
		o.meta.key = key
	}
}

func WithSecret(secret string) OptionFunc {
	return func(o *OssManager) {
		o.meta.secret = secret
	}
}
