package http

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/unique"
	"golang.org/x/net/netutil"
	"net"
	"net/http"
	"sync"
	"time"
)

func newDefaultSrvCfg() *srvCfg {
	return &srvCfg{
		logId: func() string {
			snowflake, snowflakeErr := defaultUniqueGenerator.UniqueKey(unique.Snowflake)
			if snowflakeErr != nil {
				panic(fmt.Sprintf("internal err:%v", snowflakeErr))
			}
			return snowflake
		}(),

		limit:          1e2,
		MaxHeaderBytes: 1 << 20,

		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,

		gracefully: time.Minute,
		logger:     log.DummyLogger,
		context:    context.Background(),
	}
}

func Serve(ctx context.Context, addr string, handle http.Handler, opts ...SrvOption) error {
	cfg := newDefaultSrvCfg()
	{
		for _, opt := range opts {
			opt.apply(cfg)
		}
		cfg.context = ctx
	}

	if err := parseIp(addr); err != nil {
		cfg.logger.Errorf("parseIp failed => id:%v,err:%v", cfg.logId, err)
		return err
	}

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		cfg.logger.Errorf("tcp listen:%v failed,id:%v,err:%v", addr, cfg.logId, err)
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
			cfg.logger.Errorf("server shutdown failed => id:%v,err:%v", cfg.logId, err)
			return
		}
		cfg.logger.Infof("server gracefully shutdown => id:%v", cfg.logId)

	}()

	cfg.logger.Infof("About to listen on %v,id:%v", addr, cfg.logId)
	err = srv.Serve(listener)
	wg.Wait()

	if http.ErrServerClosed == err {
		return nil
	}
	cfg.logger.Errorf("server not gracefully shutdown => id:%v,err:%v", cfg.logId, err)

	return err
}

func parseIp(ip string) error {
	switch {
	case net.ParseIP(ip).To4() == nil:
		{
			return nil
		}
	case net.ParseIP(ip).To16() == nil:
		{
			return nil
		}
	default:
		{
			return fmt.Errorf("illegal address")
		}
	}
}
