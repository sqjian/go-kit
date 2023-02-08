package validator

type OptionFunc func(*validator)

func WithValidatorType(validatorType ValidatorType) OptionFunc {
	return func(validator *validator) {
		validator.validatorType = validatorType
	}
}
