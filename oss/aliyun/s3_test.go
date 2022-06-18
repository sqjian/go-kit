package aliyun_test

import (
	"context"
	"github.com/sqjian/go-kit/oss/aliyun"
	"testing"
)

func Test_S3(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	cli, err := aliyun.NewOssCli(
		aliyun.WithKey("xxx"),
		aliyun.WithSecret("xxx"),
		aliyun.WithAddr("http://oss-cn-shanghai.aliyuncs.com"),
	)
	checkErr(err)

	var (
		bucket    = "test-oss-public"
		objectKey = "iflytek/fuck/sqjian"
		data      = []byte("sqjian")
	)
	{
		err := cli.Upload(
			context.Background(),
			bucket,
			objectKey,
			data,
		)
		checkErr(err)
	}
	{
		rst, err := cli.Download(
			context.Background(),
			bucket,
			objectKey,
		)
		checkErr(err)
		t.Log("eq:", len(data) == len(rst))
		t.Log("rst:", string(rst))
		t.Log("data:", string(data))
	}
}
