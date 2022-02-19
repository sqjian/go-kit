package json_test

import (
	"github.com/sqjian/go-kit/validator"
	"testing"
)

import _ "embed"

func TestValidateYaml(t *testing.T) {
	// Validate uses the jsonschema to validate the configuration
	//func Validate(config map[string]interface{}, version string) error {
	//	schemaData, err := _escFSByte(false, fmt.Sprintf("/data/config_schema_v%s.json", version))
	//	if err != nil {
	//	return errors.Errorf("unsupported Compose file version: %s", version)
	//}
	//
	//	schemaLoader := gojsonschema.NewStringLoader(string(schemaData))
	//	dataLoader := gojsonschema.NewGoLoader(config)
	//
	//	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	//	if err != nil {
	//	return err
	//}
	//
	//	if !result.Valid() {
	//	return toError(result)
	//}
	//
	//	return nil
	//}
	validatorInst, validatorInstErr := validator.NewValidator(validator.WithValidatorType(validator.Json))
	if validatorInstErr != nil {
		t.Fatal(validatorInstErr)
	}
	{
		t.Log(validatorInst.Validate(testSchema, testExample))
		t.Log(validatorInst.Validate(testSchema, testExampleFake))
	}
}
