package http_test

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/sqjian/go-kit/log"
	httpUtil "github.com/sqjian/go-kit/net/http"
	"net/http"
	"testing"
	"time"
)

func index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Print(w, "Welcome!\n")
}

func hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	_, _ = fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func chat(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var upgrader = websocket.Upgrader{}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("upgrade:", err)
		return
	}

	defer ws.Close()

	processMsg := func(ws *websocket.Conn) {
		_ = ws.WriteMessage(websocket.TextMessage, []byte("hello"))
	}

	processMsg(ws)
}

func TestServe(t *testing.T) {
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
		log.WithLevel("debug"),
		log.WithConsole(false),
	)

	checkErr(loggerErr)

	router := httprouter.New()
	router.GET("/", index)
	router.GET("/hello/:name", hello)
	router.GET("/chat", chat)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(time.Hour)
		cancel()
	}()

	err := httpUtil.Serve(ctx, "0.0.0.0:8888", router, httpUtil.WithServerLogger(logger))
	if err != nil {
		t.Fatal(err)
	}
}
