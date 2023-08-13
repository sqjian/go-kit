package configor

import (
	"fmt"
	vd "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/go-playground/validator/v10"
	"github.com/sqjian/go-kit/helper"
)

func validate(obj any) error {
	{
		validatorInst := validator.New()
		if validatorErr := validatorInst.Struct(obj); validatorErr != nil {
			return validatorErr
		}
	}
	{
		validatorInst := vd.New("expr").SetErrorFactory(func(failPath, msg string) error {
			return &vd.Error{
				FailPath: failPath,
				Msg:      msg,
			}
		})
		if exprErr := validatorInst.Validate(obj, true); exprErr != nil {
			return exprErr
		}
	}

	return nil
}

func LoadJsonContents(obj any, data []byte) error {
	if !helper.IsPtr(obj) {
		return fmt.Errorf("input obj should be pointer")
	}
	if len(data) == 0 {
		return fmt.Errorf("input data is empty")
	}

	viperWrapperInst := newViperWrapper()

	if initErr := viperWrapperInst.initWithJson(data); initErr != nil {
		return initErr
	}

	if unmarshalErr := viperWrapperInst.Unmarshal(obj); unmarshalErr != nil {
		return unmarshalErr
	}

	return validate(obj)
}

func LoadTomlContents(obj any, data []byte) error {
	if !helper.IsPtr(obj) {
		return fmt.Errorf("input obj should be pointer")
	}
	if len(data) == 0 {
		return fmt.Errorf("input data is empty")
	}

	viperWrapperInst := newViperWrapper()

	if initErr := viperWrapperInst.initWithToml(data); initErr != nil {
		return initErr
	}

	if unmarshalErr := viperWrapperInst.Unmarshal(obj); unmarshalErr != nil {
		return unmarshalErr
	}

	return validate(obj)
}
