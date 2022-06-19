package easyhttp_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/easylog"
	"github.com/sqjian/go-kit/net/easyhttp"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	logger, loggerErr := easylog.NewLogger(
		easylog.WithFileName("go-kit.easylog"),
		easylog.WithMaxSize(3),
		easylog.WithMaxBackups(3),
		easylog.WithMaxAge(3),
		easylog.WithLevel(easylog.Debug),
		easylog.WithConsole(false),
	)

	checkErr(loggerErr)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	rst, err := easyhttp.Do(
		context.Background(),
		easyhttp.GET,
		ts.URL,
		easyhttp.WithCliQuery(map[string]string{
			"from":    "cn",
			"to":      "en",
			"content": "你好",
		}),
		easyhttp.WithCliRetry(3),
		easyhttp.WithCliLogger(logger),
	)
	checkErr(err)
	t.Logf("rst:%v,err:%v", string(rst), err)
}

func TestDoWithId(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	logger, loggerErr := easylog.NewLogger(
		easylog.WithFileName("go-kit.easylog"),
		easylog.WithMaxSize(3),
		easylog.WithMaxBackups(3),
		easylog.WithMaxAge(3),
		easylog.WithLevel(easylog.Debug),
		easylog.WithConsole(false),
	)

	checkErr(loggerErr)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	rst, err := easyhttp.Do(
		context.Background(),
		easyhttp.GET,
		ts.URL,
		easyhttp.WithCliQuery(map[string]string{
			"from":    "cn",
			"to":      "en",
			"content": "你好",
		}),
		easyhttp.WithCliRetry(3),
		easyhttp.WithCliUniqueId("xxx"),
		easyhttp.WithCliLogger(logger),
	)
	checkErr(err)
	t.Logf("rst:%v,err:%v", string(rst), err)
}

func BenchmarkDo(b *testing.B) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	logger, loggerErr := easylog.NewLogger(
		easylog.WithFileName("go-kit.easylog"),
		easylog.WithMaxSize(3),
		easylog.WithMaxBackups(3),
		easylog.WithMaxAge(3),
		easylog.WithLevel(easylog.Debug),
		easylog.WithConsole(false),
	)

	checkErr(loggerErr)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := easyhttp.Do(
				context.Background(),
				easyhttp.GET,
				ts.URL,
				easyhttp.WithCliQuery(map[string]string{
					"from":    "cn",
					"to":      "en",
					"content": "你好",
				}),
				easyhttp.WithCliRetry(3),
				easyhttp.WithCliLogger(logger),
			)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
