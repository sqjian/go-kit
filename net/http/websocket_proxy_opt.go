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

func WithWebsocketProxyIncomeInterceptor(incomeInterceptor IncomeInterceptor) WebsocketProxyOptionFunc {
	return func(options *WebsocketProxy) {
		options.IncomeInterceptor = incomeInterceptor
	}
}

func WithWebsocketProxyOutcomeInterceptor(outcomeInterceptor OutcomeInterceptor) WebsocketProxyOptionFunc {
	return func(options *WebsocketProxy) {
		options.OutcomeInterceptor = outcomeInterceptor
	}
}
