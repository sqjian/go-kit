package es

import (
	"context"
	"testing"
)

func Test_Create(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	const (
		esAddr  = `http://xrpc.xfyun.cn:9200`
		esIndex = `go-kit`
	)

	var (
		ctx = context.Background()
	)

	cli, esCLiErr := newEsCli([]string{esAddr})
	checkErr(esCLiErr)

	exists, err := cli.indexExists(ctx, esIndex)
	checkErr(err)
	if !exists {
		err := cli.createIndex(ctx, esIndex)
		checkErr(err)
	}
}
