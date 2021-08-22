package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/retry"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

var (
	defaultHttpClient *http.Client
)

const (
	defaultDialTimeout         = 30 * time.Second
	defaultHttpTimeout         = 60 * time.Second
	defaultDialKeepAlive       = 30 * time.Second
	defaultMaxConnPerHost      = 1000
	defaultIdleConnTimeout     = time.Minute
	defaultTLSHandshakeTimeout = 2 * time.Second
)

func init() {
	defaultHttpClient = &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   defaultDialTimeout,
				KeepAlive: defaultDialKeepAlive,
			}).DialContext,
			MaxConnsPerHost:     defaultMaxConnPerHost,
			IdleConnTimeout:     defaultIdleConnTimeout,
			TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
		},
		Timeout: defaultHttpTimeout,
	}
}
func newDefaultHttpConfig() *Config {
	return &Config{
		client:  defaultHttpClient,
		context: context.Background(),
		retry:   3,
		logger:  log.DummyLogger,
	}
}
func Do(
	method Method,
	target string,
	opts ...Option,
) ([]byte, error) {
	config := newDefaultHttpConfig()

	for _, opt := range opts {
		opt.apply(config)
	}

	if err := config.context.Err(); err != nil {
		return nil, err
	}

	newReq := func() (*http.Request, error) {
		u, err := url.Parse(target)
		if err != nil {
			return nil, err
		}

		q := u.Query()
		for k, v := range config.query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		urlEncode := u.String()

		req, err := http.NewRequest(method.String(), urlEncode, bytes.NewReader(config.body))
		if err != nil {
			return nil, err
		}

		for k, v := range config.header {
			req.Header.Set(k, v)
		}
		req = req.WithContext(config.context)
		return req, nil
	}

	do := func(req *http.Request) ([]byte, error) {
		resp, err := config.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}

	traceCtx := httptrace.WithClientTrace(context.Background(), &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			config.logger.Infof("Prepare to get a connection for %s.\n", hostPort)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			config.logger.Infof("Got a connection: reused: %v, from the idle pool: %v.\n",
				info.Reused, info.WasIdle)
		},
		PutIdleConn: func(err error) {
			if err == nil {
				config.logger.Infof("Put a connection to the idle pool: ok.")
			} else {
				config.logger.Infof("Put a connection to the idle pool:", err.Error())
			}
		},
		DNSStart: func(dnsStartInfo httptrace.DNSStartInfo) {
			config.logger.Infof("Begin DNS lookup for %v.\n", dnsStartInfo.Host)
		},
		DNSDone: func(dnsDoneInfo httptrace.DNSDoneInfo) {
			config.logger.Infof("End DNS lookup, detail:%+v\n", dnsDoneInfo)
		},
		ConnectStart: func(network, addr string) {
			config.logger.Infof("Dialing... (%s:%s).\n", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			if err == nil {
				config.logger.Infof("Dial is done. (%s:%s)\n", network, addr)
			} else {
				config.logger.Infof("Dial is done with error: %s. (%s:%s)\n", err, network, addr)
			}
		},
		TLSHandshakeStart: func() {
			config.logger.Infof("Begin TLSHandshake.\n")
		},
		TLSHandshakeDone: func(connectionState tls.ConnectionState, i error) {
			config.logger.Infof("End TLSHandshake, detail:%+v\n", connectionState)
		},
		WroteHeaders: func() {
			config.logger.Infof("Wrote headers: ok.\n")
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			if info.Err == nil {
				config.logger.Infof("Wrote a request: ok.")
			} else {
				config.logger.Infof("Wrote a request:", info.Err.Error())
			}
		},
		GotFirstResponseByte: func() {
			config.logger.Infof("Got the first response byte.")
		},
	})

	var rst []byte
	err := retry.Do(
		func() error {
			req, reqErr := newReq()
			if reqErr != nil {
				return reqErr
			}
			req = req.WithContext(traceCtx)
			body, err := do(req)
			if err != nil {
				return err
			}
			rst = body
			return nil
		},
		retry.WithAttempts(uint(config.retry)),
		retry.WithContext(config.context),
	)
	return rst, err
}
