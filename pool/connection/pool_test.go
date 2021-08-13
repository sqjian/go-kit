package connection_test

import (
	"context"
	"fmt"
	"github.com/sqjian/go-kit/log/provider/zap"
	"github.com/sqjian/go-kit/pool/connection"
	"net"
	"testing"
)

func TestClientPool(t *testing.T) {
	checkErr := func(err error) {
		if err != nil {
			panic(err)
		}
	}

	logger, loggerErr := zap.NewLogger(
		zap.WithFileName("test.log"),
		zap.WithMaxSize(3),
		zap.WithMaxBackups(3),
		zap.WithMaxAge(3),
		zap.WithLevel(zap.Info),
		zap.WithConsole(false),
	)

	checkErr(loggerErr)

	pool, poolErr := connection.NewClientPool(
		context.TODO(),
		connection.WithAddress("www.baidu.com"),
		connection.WithPort("80"),
		connection.WithDial(func(ctx context.Context, address, port string) (connection interface{}, err error) {
			return net.Dial("tcp", "www.baidu.com:80")
		}),
		connection.WithClose(func(ctx context.Context, connection interface{}) (err error) {
			conn, connOk := connection.(net.Conn)
			if !connOk {
				return fmt.Errorf("can't convert to net.Conn")
			}
			return conn.Close()
		}),
		connection.WithKeepAlive(func(ctx context.Context, connection interface{}) (err error) {
			return nil
		}),
		connection.WithInitialPoolSize(1),
		connection.WithMaxPoolSize(1),
		connection.WithLogger(logger),
	)
	checkErr(poolErr)

	conn, connErr := pool.Get()
	checkErr(connErr)
	t.Log(conn.(net.Conn).RemoteAddr().String())
}
