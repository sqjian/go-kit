package pool_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/pool"
	"net"
	"testing"
)

func TestSharePool(t *testing.T) {
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
		log.WithLevel("debug"),
		log.WithConsole(false),
	)

	checkErr(loggerErr)

	const (
		address = "www.baidu.com"
		port    = "80"
	)
	p, e := pool.NewPool(
		context.TODO(),
		pool.WithType(pool.Share),
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
		pool.WithInitialPoolSize(1e3),
		pool.WithLogger(logger),
	)

	checkErr(e)

	conn, connErr := p.Get()
	checkErr(connErr)
	t.Log(conn.(net.Conn).RemoteAddr().String())
}
