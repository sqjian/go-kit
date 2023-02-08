package http_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/net/http"
	gohttp "net/http"
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
		log.WithLevel(log.Debug),
		log.WithConsole(false),
	)

	checkErr(loggerErr)

	ts := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	rst, err := http.Do(
		context.Background(),
		http.GET,
		ts.URL,
		http.WithCliQuery(map[string]string{
			"from":    "cn",
			"to":      "en",
			"content": "你好",
		}),
		http.WithCliRetry(3),
		http.WithCliLogger(logger),
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
		log.WithFileName("go-kit.easylog"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(log.Debug),
		log.WithConsole(false),
	)

	checkErr(loggerErr)

	ts := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	rst, err := http.Do(
		context.Background(),
		http.GET,
		ts.URL,
		http.WithCliQuery(map[string]string{
			"from":    "cn",
			"to":      "en",
			"content": "你好",
		}),
		http.WithCliRetry(3),
		http.WithCliUniqueId("xxx"),
		http.WithCliLogger(logger),
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
		log.WithFileName("go-kit.easylog"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel(log.Debug),
		log.WithConsole(false),
	)

	checkErr(loggerErr)

	ts := httptest.NewServer(gohttp.HandlerFunc(func(w gohttp.ResponseWriter, r *gohttp.Request) {
		time.Sleep(500 * 1e6)
		fmt.Fprintln(w, "hello")
	}))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := http.Do(
				context.Background(),
				http.GET,
				ts.URL,
				http.WithCliQuery(map[string]string{
					"from":    "cn",
					"to":      "en",
					"content": "你好",
				}),
				http.WithCliRetry(3),
				http.WithCliLogger(logger),
			)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
