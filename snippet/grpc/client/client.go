package main

import (
	retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/sqjian/go-kit/snippet/grpc/idl"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
	"time"
)

const (
	address = "localhost:50051"
)

type credential struct{}

func (c credential) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"key": "i am key",
	}, nil
}

// RequireTransportSecurity 是否开启TLS
func (c credential) RequireTransportSecurity() bool {
	return false
}

func interceptor(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	err := invoker(ctx, method, req, reply, cc, opts...)
	log.Printf("method=%s req=%v rep=%v duration=%s error=%v	\n", method, req, reply, time.Since(start), err)
	return err
}

func main() {
	retryOpts := []retry.CallOption{
		retry.WithBackoff(retry.BackoffLinear(100 * time.Millisecond)),
		retry.WithCodes(codes.NotFound, codes.Aborted),
		retry.WithMax(3),
		retry.WithPerRetryTimeout(time.Second),
	}

	var dialOpts = []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithPerRPCCredentials(new(credential)),
		grpc.WithUnaryInterceptor(interceptor),
		grpc.WithUnaryInterceptor(retry.UnaryClientInterceptor(retryOpts...)),
		grpc.WithStreamInterceptor(retry.StreamClientInterceptor(retryOpts...)),
	}

	// 连接
	conn, err := grpc.Dial(address, dialOpts...)
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// 初始化客户端
	c := idl.NewGreeterClient(conn)
	// 调用方法
	req := &idl.HelloRequest{Name: "world"}

	res, err := c.SayHello(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf(res.Message)
}
