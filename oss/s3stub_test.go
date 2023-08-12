package oss_test

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/sqjian/go-kit/oss"
)

func Example_s3stub() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}
	stub, err := oss.NewS3Stub(
		"12345678",
		"12345678",
		"http://x.x.x.x:xxxx",
	)
	checkErr(err)

	buckets, err := stub.ListBuckets()
	checkErr(err)
	spew.Dump(buckets)

	for _, bucket := range buckets {
		objs, err := stub.ListObjects(*bucket.Name)
		checkErr(err)
		spew.Printf("bucket:%v,objs:%v", *bucket.Name, objs)
	}
}
