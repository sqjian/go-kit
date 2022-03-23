package json

type Option interface {
	apply(*Validator)
}

type optionFunc func(*Validator)

func (f optionFunc) apply(log *Validator) {
	f(log)
}
