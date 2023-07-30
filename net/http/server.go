package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
	"sync"
	"time"
)

func newDefaultServerCfg() *serverConfig {
	return &serverConfig{
		limit:          1e2,
		MaxHeaderBytes: 1 << 20,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		gracefully: time.Minute,
		logger:     func() log.Log { inst, _ := log.NewLogger(log.WithLevel("dummy")); return inst }(),
		context:    context.Background(),
	}
}

func Serve(ctx context.Context, addr string, handle http.Handler, opts ...ServerOptionFunc) error {
	configInst := newDefaultServerCfg()
	{
		for _, opt := range opts {
			opt(configInst)
		}
		configInst.context = ctx
	}

	if err := parseIp(addr); err != nil {
		configInst.logger.Errorf("parseIp failed =>err:%v", err)
		return err
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		configInst.logger.Errorf("tcp listen:%v failed,err:%v", addr, err)
		return err
	}
	defer listener.Close()

	listener = netutil.LimitListener(listener, configInst.limit)
	srv := &http.Server{
		Addr:           addr,
		Handler:        handle,
		ReadTimeout:    configInst.ReadTimeout,
		WriteTimeout:   configInst.WriteTimeout,
		MaxHeaderBytes: configInst.MaxHeaderBytes,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-configInst.context.Done()
		ctx, cancel := context.WithTimeout(context.Background(), configInst.gracefully)
		defer cancel()
		if err := srv.Shutdown(ctx); nil != err {
			configInst.logger.Errorf("server shutdown failed => err:%v", err)
			return
		}
		configInst.logger.Infof("server gracefully shutdown")

	}()

	configInst.logger.Infof("About to listen on %v", addr)
	err = srv.Serve(listener)
	wg.Wait()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	configInst.logger.Errorf("server not gracefully shutdown => err:%v", err)

	return err
}

func parseIp(ip string) error {
	if net.ParseIP(ip).To4() == nil ||
		net.ParseIP(ip).To16() == nil {
		return nil
	}
	return fmt.Errorf("illegal address")
}
