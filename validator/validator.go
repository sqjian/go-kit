package validator

import (
	"github.com/sqjian/go-kit/validator/json"
)

type Validator interface {
	Validate([]byte /*schema*/, []byte /*data*/) error
}

type validator struct {
	validatorType ValidatorType
}

func newDefaultValidatorConfig() *validator {
	return &validator{}
}

func NewValidator(opts ...OptionFunc) (Validator, error) {

	configInst := newDefaultValidatorConfig()

	for _, opt := range opts {
		opt(configInst)
	}

	switch configInst.validatorType {
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
