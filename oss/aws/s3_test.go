package aws_test

import (
	"context"
	"github.com/sqjian/go-kit/oss/aws"
	"testing"
)

func Test_S3(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	cli, err := aws.NewS3Cli(
		aws.WithKey("12345678"),
		aws.WithSecret("12345678"),
		aws.WithAddr("http://172.31.243.215:9091"),
		aws.WithConcurrency(3),
	)
	checkErr(err)

	var (
		bucket    = "sqjian"
		objectKey = "sqjian"
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
