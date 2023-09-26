package configor_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sqjian/go-kit/configor"
	"testing"
)

func TestLoadJsonContents(t *testing.T) {
	obj := &struct {
		Salary int    `validate:"$>5000" json:"salary,omitempty"`
		Tel    string `validate:"regexp('^(13[0-9]|14[01456879]|15[0-35-9]|16[2567]|17[0-8]|18[0-9]|19[0-35-9])\\d{8}$')" json:"tel,omitempty"`
	}{}
	data := []byte(`{"name": "sqjian","age": "12","gender": "male","salary": 9000,"tel": "19909697214"}`)

	spew.Dump(configor.LoadJsonContents(obj, data))
}

func TestLoadNestJsonContents(t *testing.T) {
	obj := &struct {
		OssConfig struct {
			S3AccessKey string `json:"accesskey" validate:"len($)>1"`
			S3SecretKey string `json:"secretkey" validate:"len($)>1"`
		} `json:"oss"`
		RdbConfig struct {
			DbName    string `json:"db" validate:"len($)>1"`
			TableName string `json:"table" validate:"len($)>1"`
		} `json:"rdb"`
	}{}
	data := []byte(`{
    "rdb": {
        "db": "xxx",
        "table": "xxx"
    },
    "oss": {
        "accesskey": "xxx",
        "secretkey": "xxx"
    }
}
`)
	spew.Dump(configor.LoadJsonContents(obj, data))
}
