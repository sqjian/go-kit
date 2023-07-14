package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/retry"
	"github.com/sqjian/go-kit/uid"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

var (
	defaultHttpClient      *http.Client
	defaultUniqueGenerator uid.Uid
)

const (
	defaultMaxConnPerHost = 1000

	defaultDialTimeout         = 30 * time.Second
	defaultHttpTimeout         = 60 * time.Second
	defaultDialKeepAlive       = 30 * time.Second
	defaultIdleConnTimeout     = time.Minute
	defaultTLSHandshakeTimeout = 2 * time.Second
)

func GetDefaultHttpClient() *http.Client {
	return defaultHttpClient
}

func init() {
	defaultHttpClient = &http.Client{
		Transport: &http.Transport{
			Proxy:               http.ProxyFromEnvironment,
			DialContext:         (&net.Dialer{Timeout: defaultDialTimeout, KeepAlive: defaultDialKeepAlive}).DialContext,
			MaxConnsPerHost:     defaultMaxConnPerHost,
			IdleConnTimeout:     defaultIdleConnTimeout,
			TLSHandshakeTimeout: defaultTLSHandshakeTimeout,
		},
		Timeout: defaultHttpTimeout,
	}

	defaultUniqueGenerator = func() uid.Uid {
		uniqueGenerator, uniqueGeneratorErr := uid.NewGenerator(
			uid.Snowflake,
			uid.WithSnowflakeNodeId(1),
		)
		if uniqueGeneratorErr != nil {
			panic(fmt.Sprintf("internal err:%v", uniqueGeneratorErr))
		}
		return uniqueGenerator
	}()
}

func newDefaultCliCfg() *clientConfig {
	return &clientConfig{
		retry:   3,
		trace:   true,
		Log:     func() log.Log { inst, _ := log.NewLogger(log.WithLevel(log.Dummy)); return inst }(),
		client:  defaultHttpClient,
		context: context.Background(),
	}
}

func genTraceCtx(config *clientConfig) context.Context {
	traceCtx := httptrace.WithClientTrace(config.context, &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			config.Infof("Prepare to get a connection for %s.", hostPort)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			config.Infof("Got a connection: reused: %v, from the idle pool: %v.", info.Reused, info.WasIdle)
		},
		PutIdleConn: func(err error) {
			if err == nil {
				config.Infof("Put a connection to the idle pool: ok.")
			} else {
				config.Infof("Put a connection to the idle pool:%v", err.Error())
			}
		},
		DNSStart: func(dnsStartInfo httptrace.DNSStartInfo) {
			config.Infof("Begin DNS lookup for %v.", dnsStartInfo.Host)
		},
		DNSDone: func(dnsDoneInfo httptrace.DNSDoneInfo) {
			config.Infof("End DNS lookup, detail:%+v\n", dnsDoneInfo)
		},
		ConnectStart: func(network, addr string) {
			config.Infof("Dialing... (%s:%s).", network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			if err == nil {
				config.Infof("Dial is done. (%s:%s)", network, addr)
			} else {
				config.Infof("Dial is done with error: %s. (%s:%s)", err, network, addr)
			}
		},
		TLSHandshakeStart: func() {
			config.Infof("Begin TLSHandshake.")
		},
		TLSHandshakeDone: func(connectionState tls.ConnectionState, i error) {
			config.Infof("End TLSHandshake, detail:%+v", connectionState)
		},
		WroteHeaders: func() {
			config.Infof("Wrote headers: ok.")
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			if info.Err == nil {
				config.Infof("Wrote a request: ok.")
			} else {
				config.Infof("Wrote a request:%v", info.Err.Error())
			}
		},
		GotFirstResponseByte: func() {
			config.Infof("Got the first response byte.")
		},
	})
	return traceCtx
}

func genReq(method Method, target string, cfg *clientConfig) (*http.Request, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}

	q := u.Query()
	for k, v := range cfg.query {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()
	urlEncode := u.String()

	req, err := http.NewRequest(
		method.String(),
		urlEncode,
		bytes.NewReader(cfg.body),
	)
	if err != nil {
		return nil, err
	}

	for k, v := range cfg.header {
		req.Header.Set(k, v)
	}

	if cfg.trace {
		req = req.WithContext(genTraceCtx(cfg))
	} else {
		req = req.WithContext(cfg.context)
	}

	return req, nil
}

func Do(ctx context.Context, method Method, target string, opts ...CliOptionFunc) ([]byte, error) {
	cfg := newDefaultCliCfg()
	cfg.context = ctx

	for _, opt := range opts {
		opt(cfg)
	}

	if err := cfg.context.Err(); err != nil {
		cfg.Errorf("context.Err not nil =>err:%v", err)
		return nil, err
	}

	do := func(req *http.Request) ([]byte, error) {
		resp, err := cfg.client.Do(req)
		if err != nil {
			cfg.Errorf("client.do failed =>err:%v", err)
			return nil, err
		}
		defer resp.Body.Close()

		respData, respDataErr := io.ReadAll(resp.Body)
		if respDataErr != nil {
			return nil, respDataErr
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("statusCode:%v not equal 200,raw:%v", resp.StatusCode, string(respData))
		}
		return respData, respDataErr
	}

	var rst []byte
	err := retry.Do(
		func() error {
			req, reqErr := genReq(method, target, cfg)
			if reqErr != nil {
				cfg.Errorf("http.Do->genReq failed,err:%v", reqErr)
				return reqErr
			}
			body, err := do(req)
			if err != nil {
				cfg.Errorf("http.Do->do failed,err:%v", err)
				return err
			}
			rst = body
			return nil
		},
		retry.WithAttempts(uint(cfg.retry)),
		retry.WithContext(cfg.context),
	)
	return rst, err
}
