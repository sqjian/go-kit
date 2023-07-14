package http

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
	"sync"
	"time"
)

func newDefaultSrvCfg() *serverConfig {
	return &serverConfig{
		limit:          1e2,
		MaxHeaderBytes: 1 << 20,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		gracefully: time.Minute,
		logger:     func() log.Log { inst, _ := log.NewLogger(log.WithLevel(log.Dummy)); return inst }(),
		context:    context.Background(),
	}
}

func Serve(ctx context.Context, addr string, handle http.Handler, opts ...ServerOptionFunc) error {
	cfg := newDefaultSrvCfg()
	{
		for _, opt := range opts {
			opt(cfg)
		}
		cfg.context = ctx
	}

	if err := parseIp(addr); err != nil {
		cfg.logger.Errorf("parseIp failed =>err:%v", err)
		return err
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		cfg.logger.Errorf("tcp listen:%v failed,err:%v", addr, err)
		return err
	}
	defer listener.Close()

	listener = netutil.LimitListener(listener, cfg.limit)
	srv := &http.Server{
		Addr:           addr,
		Handler:        handle,
		ReadTimeout:    cfg.ReadTimeout,
		WriteTimeout:   cfg.WriteTimeout,
		MaxHeaderBytes: cfg.MaxHeaderBytes,
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		<-cfg.context.Done()
		ctx, cancel := context.WithTimeout(context.Background(), cfg.gracefully)
		defer cancel()
		if err := srv.Shutdown(ctx); nil != err {
			cfg.logger.Errorf("server shutdown failed => err:%v", err)
			return
		}
		cfg.logger.Infof("server gracefully shutdown")

	}()

	cfg.logger.Infof("About to listen on %v", addr)
	err = srv.Serve(listener)
	wg.Wait()

	if http.ErrServerClosed == err {
		return nil
	}
	cfg.logger.Errorf("server not gracefully shutdown => err:%v", err)

	return err
}

func parseIp(ip string) error {
	if net.ParseIP(ip).To4() == nil ||
		net.ParseIP(ip).To16() == nil {
		return nil
	}
	return fmt.Errorf("illegal address")
}
