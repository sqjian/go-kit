package easyhttp_test

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sqjian/go-kit/easylog"
	httpUtil "github.com/sqjian/go-kit/net/easyhttp"
	"net/http"
	"testing"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Print(w, "Welcome!\n")
}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func TestServe(t *testing.T) {
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

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 3)
		cancel()
	}()

	err := httpUtil.Serve(ctx, "0.0.0.0:8888", router, httpUtil.WithSrvLogger(logger))
	if err != nil {
		t.Fatal(err)
	}
}
