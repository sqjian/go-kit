package aes

type Option interface {
	apply(*aes)
}

type optionFunc func(*aes)

func (f optionFunc) apply(a *aes) {
	f(a)
}

func WithAesMode(mode Mode) Option {
	return optionFunc(func(aes *aes) {
		aes.mode = mode
	})
}
