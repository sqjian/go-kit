package digest

type Option interface {
	apply(*generator)
}

type optionFunc func(*generator)

func (f optionFunc) apply(log *generator) {
	f(log)
}
