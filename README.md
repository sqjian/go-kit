## service toolkit

[![Go](https://github.com/sqjian/go-kit/actions/workflows/go-kit.yml/badge.svg)](https://github.com/sqjian/go-kit/actions/workflows/go-kit.yml)
[![GoDoc](https://godoc.org/github.com/sqjian/go-kit?status.svg&style=flat-square)](http://godoc.org/github.com/sqjian/go-kit)
[![Go Report Card](https://goreportcard.com/badge/github.com/sqjian/go-kit?style=flat-square)](https://goreportcard.com/report/github.com/sqjian/go-kit)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/sqjian/go-kit)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/sqjian/go-kit)

## TODO

- configurator support nest struct
```
type Cfg struct {
	OssConfig struct {
		S3AccessKey string `validate:"required" json:"accesskey"`
		S3SecretKey string `json:"secretkey" validate:"required"`
		S3Host      string `json:"host" validate:"required"`
		S3Timeout   int64  `json:"timeout" validate:"required"`
		S3Bucket    string `json:"buckets" validate:"required"`
		S3Debug     bool   `json:"debug"`
	} `json:"oss"`
	RdbConfig struct {
		RdbAddr      string `json:"addr" validate:"required"`
		RdbPort      string `json:"port" validate:"required"`
		UserName     string `json:"user" validate:"required"`
		Password     string `json:"passwd" validate:"required"`
		MaxLifeTime  int64  `json:"maxlifetime" validate:"required"`
		MaxIdleConns int64  `json:"maxidleconns" validate:"required"`
		DbName       string `json:"db" validate:"required"`
		TableName    string `json:"table" validate:"required"`
	} `json:"rdb"`
}

{
  "rdb": {
    "addr": "xxx",
    "port": "xxx",
    "db": "xxx",
    "maxidleconns": 1,
    "maxlifetime": 1,
    "passwd": "xxx",
    "table": "xxx",
    "user": "xxx"
  },
  "oss": {
    "host": "xxx",
    "accesskey": "xxx",
    "secretkey": "xxx",
    "buckets": "xxx"
  }
}
```