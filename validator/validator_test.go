package validator_test

import (
	_ "embed"
)

////go:embed test/test-schema.json
//var testSchema []byte
//
////go:embed test/test-example.json
//var testExample []byte
//
////go:embed test/test-example-fake.json
//var testExampleFake []byte

// func TestValidateJson(t *testing.T) {
// 	checkErr := func(err error) {
// 		if err != nil {
// 			t.Log(err.Error())
// 			os.Exit(0)
// 		}
// 	}
// 	validatorInst, validatorInstErr := validator.NewValidator(validator.Json)
// 	checkErr(validatorInstErr)
// 	{
// 		t.Log(validatorInst.Validate(testSchema, testExample))
// 		t.Log(validatorInst.Validate(testSchema, testExampleFake) != nil)
// 	}
// }
