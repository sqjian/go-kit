package main

import (
	"context"
	"github.com/sqjian/go-kit/snippet/grpc/idl"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

const (
	address = "localhost:50051"
)

type server struct{}

func (s *server) SayHello(ctx context.Context, in *idl.HelloRequest) (*idl.HelloReply, error) {
	return &idl.HelloReply{Message: "Hello " + in.Name}, nil
}

func auth(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "无Token认证信息")
	}
	var (
		key string
	)
	if val, ok := md["key"]; ok {
		key = val[0]
	}
	if key != "i am key" {
		return status.Errorf(codes.Unauthenticated, "Token认证信息无效: key =%s", key)
	}
	return nil
}

func interceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	err := auth(ctx)
	if err != nil {
		return nil, err
	}
	return handler(ctx, req)
}

func main() {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	opts = append(opts, grpc.UnaryInterceptor(interceptor))

	s := grpc.NewServer(opts...)
	idl.RegisterGreeterServer(s, &server{})
	reflection.Register(s)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
