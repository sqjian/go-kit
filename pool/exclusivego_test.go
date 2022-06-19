package pool_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/easylog"
	"github.com/sqjian/go-kit/pool"
	"net"
	"testing"
)

func TestExclusivePool(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	logger, loggerErr := easylog.NewLogger(
		easylog.WithFileName("go-kit.easylog"),
		easylog.WithMaxSize(3),
		easylog.WithMaxBackups(3),
		easylog.WithMaxAge(3),
		easylog.WithLevel(easylog.Warn),
		easylog.WithConsole(false),
	)

	checkErr(loggerErr)

	const (
		address = "www.baidu.com"
		port    = "80"
	)
	pool, poolErr := pool.NewPool(
		context.TODO(),
		pool.WithType(pool.Exclusive),
		pool.WithAddress(address),
		pool.WithPort(port),
		pool.WithDial(func(ctx context.Context, address, port string) (connection interface{}, err error) {
			return net.Dial("tcp", fmt.Sprintf("%v:%v", address, port))
		}),
		pool.WithClose(func(ctx context.Context, connection interface{}) (err error) {
			conn, connOk := connection.(net.Conn)
			if !connOk {
				return fmt.Errorf("can't convert to net.Conn")
			}
			return conn.Close()
		}),
		pool.WithKeepAlive(func(ctx context.Context, connection interface{}) (err error) {
			return nil
		}),
		pool.WithInitialPoolSize(1),
		pool.WithMaxPoolSize(1),
		pool.WithLogger(logger),
	)

	checkErr(poolErr)

	conn, connErr := pool.Get()
	checkErr(connErr)
	t.Log(conn.(net.Conn).RemoteAddr().String())
}
