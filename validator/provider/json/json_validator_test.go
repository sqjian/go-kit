package json_test

import (
	"github.com/sqjian/go-kit/validator"
	"testing"
)

import _ "embed"

//go:embed testdata/test-schema.json
var testSchema []byte

//go:embed testdata/test-example.json
var testExample []byte

//go:embed testdata/test-example-fake.json
var testExampleFake []byte

func TestValidateJson(t *testing.T) {
	validatorInst, _ := validator.NewValidator()
	{
		t.Log(validatorInst.Validate(testSchema, testExample))
		t.Log(validatorInst.Validate(testSchema, testExampleFake))
	}
}
