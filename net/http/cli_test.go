package http_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/log/vars"
	httpUtil "github.com/sqjian/go-kit/net/http"
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

	logger, loggerErr := log.NewLogger(
		log.WithFileName("go-kit.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(vars.Debug),
		log.WithLogType(vars.Zap),
		log.WithConsole(false),
	)

	checkErr(loggerErr)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	rst, err := httpUtil.Do(
		context.Background(),
		httpUtil.GET,
		ts.URL,
		httpUtil.WithCliQuery(map[string]string{
			"from":    "cn",
			"to":      "en",
			"content": "你好",
		}),
		httpUtil.WithCliRetry(3),
		httpUtil.WithCliLogger(logger),
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

	logger, loggerErr := log.NewLogger(
		log.WithFileName("go-kit.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(vars.Debug),
		log.WithLogType(vars.Zap),
		log.WithConsole(false),
	)

	checkErr(loggerErr)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	rst, err := httpUtil.Do(
		context.Background(),
		httpUtil.GET,
		ts.URL,
		httpUtil.WithCliQuery(map[string]string{
			"from":    "cn",
			"to":      "en",
			"content": "你好",
		}),
		httpUtil.WithCliRetry(3),
		httpUtil.WithCliUniqueId("xxx"),
		httpUtil.WithCliLogger(logger),
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

	logger, loggerErr := log.NewLogger(
		log.WithFileName("go-kit.log"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(vars.Debug),
		log.WithLogType(vars.Zap),
		log.WithConsole(false),
	)

	checkErr(loggerErr)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := httpUtil.Do(
				context.Background(),
				httpUtil.GET,
				ts.URL,
				httpUtil.WithCliQuery(map[string]string{
					"from":    "cn",
					"to":      "en",
					"content": "你好",
				}),
				httpUtil.WithCliRetry(3),
				httpUtil.WithCliLogger(logger),
			)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
