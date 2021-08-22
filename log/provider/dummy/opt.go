package dummy

type Option interface {
	apply(*logger)
}

type optionFunc func(*logger)

func (f optionFunc) apply(log *logger) {
	f(log)
}
