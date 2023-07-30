package http

import (
	"github.com/sqjian/go-kit/log"
)

type WebsocketProxyOptionFunc func(*WebsocketProxy)

func WithWebsocketProxyLogger(logger log.Log) WebsocketProxyOptionFunc {
	return func(options *WebsocketProxy) {
		options.Logger = logger
	}
}

func WithWebsocketProxyInterceptor(interceptor Interceptor) WebsocketProxyOptionFunc {
	return func(options *WebsocketProxy) {
		options.Interceptors = append(options.Interceptors, interceptor)
	}
}

//func Process(messageType int, data []byte) (int, []byte) {
//
//}
