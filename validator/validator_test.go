package validator_test

import (
	_ "embed"
	"github.com/sqjian/go-kit/validator"
	"os"
	"testing"
)

//go:embed testdata/test-schema.json
var testSchema []byte

//go:embed testdata/test-example.json
var testExample []byte

//go:embed testdata/test-example-fake.json
var testExampleFake []byte

func TestValidateJson(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Log(err.Error())
			os.Exit(0)
		}
	}
	validatorInst, validatorInstErr := validator.NewValidator(
		validator.WithValidatorType(validator.Json),
	)
	checkErr(validatorInstErr)
	{
		t.Log(validatorInst.Validate(testSchema, testExample))
		t.Log(validatorInst.Validate(testSchema, testExampleFake))
	}
}