package configor

import (
	"encoding/json"
	"fmt"
	expr "github.com/bytedance/go-tagexpr/v2/validator"
	. "github.com/sqjian/go-kit/encoding/json"
	"github.com/sqjian/go-kit/helper"
)

func validate(obj any) error {
	validatorInst := expr.New("validate").SetErrorFactory(func(failPath, msg string) error {
		return &expr.Error{
			FailPath: failPath,
			Msg:      msg,
		}
	})
	return validatorInst.Validate(obj, true)
}

func LoadJsonContents(obj any, data []byte) error {
	if !helper.IsPtr(obj) {
		return fmt.Errorf("input obj should be pointer")
	}
	if len(data) == 0 {
		return fmt.Errorf("input data is empty")
	}

	if unmarshalErr := json.Unmarshal(Standardize(data), obj); unmarshalErr != nil {
		return unmarshalErr
	}

	return validate(obj)
}
