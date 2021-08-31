package validator

type Option interface {
	apply(*validator)
}

type optionFunc func(*validator)

func (f optionFunc) apply(validator *validator) {
	f(validator)
}

func WithValidatorType(validatorType ValidatorType) Option {
	return optionFunc(func(validator *validator) {
		validator.validatorType = validatorType
	})
}
