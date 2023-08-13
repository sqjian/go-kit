package configor

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/sqjian/go-kit/helper"
)

func checkStruct(obj any) error {
	return validator.New().Struct(obj)
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

	return checkStruct(obj)
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

	return checkStruct(obj)
}
