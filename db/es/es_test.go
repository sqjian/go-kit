package es

import (
	"context"
	"github.com/sqjian/go-kit/net/http"
	"testing"
)

func Test_Es(t *testing.T) {
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

	cli, esCLiErr := newEsCli(
		WithHosts(esAddr),
		WithDebugInfo(true),
		WithHttpClient(http.GetDefaultHttpClient()),
	)
	checkErr(esCLiErr)

	t.Log("1")
	exists, err := cli.indexExists(ctx, esIndex)
	checkErr(err)
	if !exists {
		t.Log("2")
		err := cli.createIndex(ctx, esIndex)
		checkErr(err)
	}
}
