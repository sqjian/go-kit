package validator

import "github.com/sqjian/go-kit/validator/provider/json"

type Validator interface {
	Validate([]byte /*schema*/, []byte /*data*/) error
}

type validator struct {
	validatorType ValidatorType
}

func newDefaultValidatorConfig() *validator {
	return &validator{}
}

func NewValidator(opts ...Option) (Validator, error) {

	validatorInst := newDefaultValidatorConfig()

	for _, opt := range opts {
		opt.apply(validatorInst)
	}

	switch validatorInst.validatorType {
	case Json:
		{
			return json.NewValidator()
		}
	default:
		{
			return nil, ErrWrapper(IllegalKeyType)
		}
	}
}
