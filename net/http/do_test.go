package http_test

import (
	"context"
	"fmt"
	httpUtil "github.com/sqjian/go-kit/net/http"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
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

const (
	testPort = 8888
)

func TestServer(t *testing.T) {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = fmt.Fprintf(writer, "Hello, there\n")
	})
	fmt.Println("Server started at port " + strconv.Itoa(testPort))
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(testPort), nil))
}

func BenchmarkDo(b *testing.B) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := httpUtil.Do(
				"http://172.31.243.215:"+strconv.Itoa(testPort),
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
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}
