package aliyun_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/oss/aliyun"
)

func Example() {
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
		fmt.Println("eq:", len(data) == len(rst))
		fmt.Println("rst:", string(rst))
		fmt.Println("data:", string(data))
	}
}
