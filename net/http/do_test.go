package http_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/log/preset"
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
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	rst, err := httpUtil.Do(
		ts.URL,
		httpUtil.GET,
		map[string]string{
			"from":    "cn",
			"to":      "en",
			"content": "你好",
		},
		nil,
		nil,
		httpUtil.WithContext(func() context.Context {
			ctx, _ := context.WithTimeout(context.Background(), 1000*time.Millisecond)
			return ctx
		}()),
		httpUtil.WithRetry(3),
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
		log.WithLevel(preset.Debug),
		log.WithLogType(preset.Zap),
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
				ts.URL,
				httpUtil.GET,
				map[string]string{
					"from":    "cn",
					"to":      "en",
					"content": "你好",
				},
				nil,
				nil,
				httpUtil.WithContext(func() context.Context {
					ctx, _ := context.WithTimeout(context.Background(), 1000*time.Millisecond)
					return ctx
				}()),
				httpUtil.WithRetry(3),
				httpUtil.WithLogger(logger),
			)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
