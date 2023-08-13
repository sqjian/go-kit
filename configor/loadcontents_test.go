package configor_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sqjian/go-kit/configor"
	"testing"
)

func TestLoadJsonContents(t *testing.T) {
	obj := &struct {
		Name   string `validate:"required" json:"name,omitempty"`
		Age    int    `validate:"gte=10,lte=130" json:"age,omitempty"`
		Gender string `validate:"oneof=male female" json:"gender,omitempty"`
		Salary int    `expr:"$>5000" json:"salary,omitempty"`
		Tel    string `expr:"regexp('^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$')" json:"tel,omitempty"`
	}{}
	data := []byte(`{"name": "sqjian","age": 12,"gender": "male","salary": 9000,"tel": 19909697214}`)

	spew.Dump(configor.LoadJsonContents(obj, data))
}

func TestLoadTomlContents(t *testing.T) {
	obj := &struct {
		Name   string `validate:"required" json:"name,omitempty"`
		Age    int    `validate:"gte=10,lte=130" json:"age,omitempty"`
		Gender string `validate:"oneof=male female" json:"gender,omitempty"`
		Salary int    `expr:"$>5000" json:"salary,omitempty"`
		Tel    string `expr:"regexp('^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$')" json:"tel,omitempty"`
	}{}
	data := []byte(`name = 'sqjian'
age = 90
gender = 'male'
salary = 4000
tel = 19909697214`)

	spew.Dump(configor.LoadTomlContents(obj, data))
}
