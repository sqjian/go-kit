package s3_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/sqjian/go-kit/oss/s3"
	"os"
)

var awsConfig *aws.Config

func init() {
	cfg, err := s3.NewAwsConfig(
		s3.WithAccessKey("xxx"),
		s3.WithSecretKey("xxx"),
		s3.WithEndpoint("http://x.x.x.x:xxxx"),
	)
	if err != nil {
		panic(err)
	}
	awsConfig = cfg
}

func Example_bucketFiles() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	cli, err := s3.NewS3Cli(
		s3.WithAwsConfig(awsConfig),
		s3.WithProgressOutput(os.Stderr),
	)
	checkErr(err)

	var (
		bucket     = "xxx"
		prefixHint = ""
	)
	{
		files, err := cli.BucketFiles(
			context.Background(),
			bucket,
			prefixHint,
		)
		checkErr(err)
		fmt.Println(files)
	}
}

func Example_uploadFile() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	cli, err := s3.NewS3Cli(
		s3.WithAwsConfig(awsConfig),
		s3.WithProgressOutput(os.Stderr),
	)
	checkErr(err)

	var (
		bucket    = "xxx"
		objectKey = "xxx"
		data      = bytes.Repeat([]byte("x"), 1024*1024*1024)
	)
	{
		err := cli.UploadFile(
			context.Background(),
			bucket,
			objectKey,
			s3.NewFsFile("", data, 0644),
		)
		checkErr(err)
	}
}

func Example_downloadFile() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	cli, err := s3.NewS3Cli(
		s3.WithAwsConfig(awsConfig),
		s3.WithProgressOutput(os.Stderr),
	)
	checkErr(err)

	var (
		bucket    = "xxx"
		objectKey = "xxx"
		buf       = manager.NewWriteAtBuffer(nil)
	)
	{

		err := cli.DownloadFile(
			context.Background(),
			bucket,
			objectKey,
			buf,
		)
		checkErr(err)
		fmt.Println(len(buf.Bytes()))
	}
}

func Example_deleteFile() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	cli, err := s3.NewS3Cli(
		s3.WithAwsConfig(awsConfig),
		s3.WithProgressOutput(os.Stdin),
	)
	checkErr(err)

	var (
		bucket    = "xxx"
		objectKey = "xxx"
	)
	{

		err := cli.DeleteFile(
			context.Background(),
			bucket,
			objectKey,
		)
		checkErr(err)
	}
}
