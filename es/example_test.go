package es_test

import (
	"context"
	"github.com/sqjian/go-kit/es"
	"github.com/sqjian/go-kit/net/http"
	"log"
	"strings"
)

func Example() {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	const (
		esAddr  = `https://x.x.x.x:xxxx`
		esIndex = `go-kit`
	)

	var (
		ctx = context.Background()
	)

	cli, esCLiErr := es.NewEsCli(
		es.WithHosts(esAddr),
		es.WithDebugInfo(true),
		es.WithHttpClient(http.GetDefaultHttpClient()),
	)
	checkErr(esCLiErr)

	{
		exists, err := cli.IndexExists(ctx, esIndex)
		checkErr(err)
		if !exists {
			err := cli.CreateIndex(ctx, esIndex)
			checkErr(err)
		}
	}
	{
		err := cli.WriteDocs(
			ctx,
			esIndex,
			map[string]any{
				"name": "sqjian",
				"age":  18,
			},
		)
		checkErr(err)
	}
	{
		rst, err := cli.QueryDocs(
			ctx,
			esIndex,
			map[string]any{
				"name": "sqjian",
				"age":  18,
			},
		)
		checkErr(err)
		log.Printf("rst:%v", strings.Join(rst, "\n"))
	}
}
