package http_test

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/log/vars"
	httpUtil "github.com/sqjian/go-kit/net/http"
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

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Second * 3)
		cancel()
	}()

	err := httpUtil.Serve(ctx, "0.0.0.0:80", router, httpUtil.WithSrvLogger(logger))
	if err != nil {
		t.Fatal(err)
	}
}
