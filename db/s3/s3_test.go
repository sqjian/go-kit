package s3

import (
	"context"
	"testing"
)

func Test_S3(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	cli, err := newS3Cli(
		WithKey("12345678"),
		WithSecret("12345678"),
		WithAddr("http://172.31.243.215:9091"),
		WithConcurrency(3),
	)
	checkErr(err)

	var (
		bucket    = "sqjian"
		objectKey = "sqjian"
		data      = []byte("sqjian")
	)

	{
		err := cli.upload(
			context.Background(),
			bucket,
			objectKey,
			data,
		)
		checkErr(err)
	}
	{
		rst, err := cli.download(
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
