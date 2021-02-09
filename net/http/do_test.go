package http_test

import (
	"context"
	"fmt"
	httpUtil "github.com/sqjian/go-kit/net/http"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Logf("ts.input:%v", r.URL.String())
		fmt.Fprintln(w, "hello")
	}))
	defer ts.Close()

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
	t.Logf("rst:%v,err:%v", string(rst), err)
}
