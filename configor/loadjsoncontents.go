package configor

import (
	"encoding/json"
	"fmt"
	expr "github.com/bytedance/go-tagexpr/v2/validator"
	"github.com/davecgh/go-spew/spew"
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

	standardizeData, standardizeDataErr := Standardize(data)
	if standardizeDataErr != nil {
		return standardizeDataErr
	}

	if unmarshalErr := json.Unmarshal(standardizeData, obj); unmarshalErr != nil {
		return unmarshalErr
	}

	spew.Dump(obj)

	return validate(obj)
}
