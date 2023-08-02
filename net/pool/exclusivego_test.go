package pool_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	pool2 "github.com/sqjian/go-kit/net/pool"
	"net"
	"testing"
)

func TestExclusivePool(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	logger, loggerErr := log.NewLogger(
		log.WithFileName("go-kit.easylog"),
		log.WithMaxSize(3),
		log.WithMaxBackups(3),
		log.WithMaxAge(3),
		log.WithLevel("warn"),
		log.WithConsole(false),
	)

	checkErr(loggerErr)

	const (
		address = "www.baidu.com"
		port    = "80"
	)
	p, e := pool2.NewPool(
		context.TODO(),
		pool2.WithType(pool2.Exclusive),
		pool2.WithAddress(address),
		pool2.WithPort(port),
		pool2.WithDial(func(ctx context.Context, address, port string) (connection interface{}, err error) {
			return net.Dial("tcp", fmt.Sprintf("%v:%v", address, port))
		}),
		pool2.WithClose(func(ctx context.Context, connection interface{}) (err error) {
			conn, connOk := connection.(net.Conn)
			if !connOk {
				return fmt.Errorf("can't convert to net.Conn")
			}
			return conn.Close()
		}),
		pool2.WithKeepAlive(func(ctx context.Context, connection interface{}) (err error) {
			return nil
		}),
		pool2.WithInitialPoolSize(1),
		pool2.WithMaxPoolSize(1),
		pool2.WithLogger(logger),
	)

	checkErr(e)

	conn, connErr := p.Get()
	checkErr(connErr)
	t.Log(conn.(net.Conn).RemoteAddr().String())
}
