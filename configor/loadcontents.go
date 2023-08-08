package configor

import "github.com/go-playground/validator/v10"

func LoadJsonContents(obj any, data []byte) error {
	viperWrapperInst := newViperWrapper()
	if initErr := viperWrapperInst.initWithJson(data); initErr != nil {
		return initErr
	}
	if unmarshalErr := viperWrapperInst.Unmarshal(obj); unmarshalErr != nil {
		return unmarshalErr
	}
	if structErr := validator.New().Struct(obj); structErr != nil {
		return structErr
	}
	return nil
}
