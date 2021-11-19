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
)

type Middleware struct {
	next    http.Handler
	message string
}

func NewMiddleware(next http.Handler, message string) *Middleware {
	return &Middleware{next: next, message: message}
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("msg: %s, Method: %s, URI: %s\n", m.message, r.Method, r.RequestURI)
	m.next.ServeHTTP(w, r)
}

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
	middleRouter := NewMiddleware(router, "I'm a transitional middleware")
	finalRouter := NewMiddleware(middleRouter, "I'm final middleware")

	err := httpUtil.Serve(context.Background(), "0.0.0.0:80", finalRouter, httpUtil.WithSrvLogger(logger))
	if err != nil {
		t.Fatal(err)
	}
}
