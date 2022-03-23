package es

import (
	"context"
	"github.com/sqjian/go-kit/net/http"
	"strings"
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

	{
		exists, err := cli.indexExists(ctx, esIndex)
		checkErr(err)
		if !exists {
			err := cli.createIndex(ctx, esIndex)
			checkErr(err)
		}
	}
	{
		err := cli.writeDocs(
			ctx,
			esIndex,
			map[string]interface{}{
				"name": "sqjian",
				"age":  18,
			},
		)
		checkErr(err)
	}
	{
		rst, err := cli.queryDocs(
			ctx,
			esIndex,
			map[string]interface{}{
				"name": "sqjian",
				"age":  18,
			},
		)
		checkErr(err)
		t.Logf("rst:%v", strings.Join(rst, "\n"))
	}
}
