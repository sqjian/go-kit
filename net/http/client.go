package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/sqjian/go-kit/log"
	"github.com/sqjian/go-kit/retry"
	"github.com/sqjian/go-kit/unique"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"time"
)

var (
	defaultHttpClient      *http.Client
	defaultUniqueGenerator unique.Generator
)

const (
	defaultDialTimeout         = 30 * time.Second
	defaultHttpTimeout         = 60 * time.Second
	defaultDialKeepAlive       = 30 * time.Second
	defaultMaxConnPerHost      = 1000
	defaultIdleConnTimeout     = time.Minute
	defaultTLSHandshakeTimeout = 2 * time.Second
)

const (
	defaultBodyVerbose = 2048
)

func GetDefaultHttpClient() *http.Client {
	return defaultHttpClient
}

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
	defaultUniqueGenerator = func() unique.Generator {
		uniqueGenerator, uniqueGeneratorErr := unique.NewGenerator(
			unique.WithSnowflakeNodeId(1),
		)
		if uniqueGeneratorErr != nil {
			panic(fmt.Sprintf("internal err:%v", uniqueGeneratorErr))
		}
		return uniqueGenerator
	}()
}

func newDefaultHttpConfig() *Config {
	return &Config{
		retry:   3,
		trace:   true,
		logger:  log.DummyLogger,
		client:  defaultHttpClient,
		context: context.Background(),
		logId: func() string {
			snowflake, snowflakeErr := defaultUniqueGenerator.UniqueKey(unique.Snowflake)
			if snowflakeErr != nil {
				panic(fmt.Sprintf("internal err:%v", snowflakeErr))
			}
			return snowflake
		}(),
	}
}

func genTraceCtx(config *Config) context.Context {
	traceCtx := httptrace.WithClientTrace(config.context, &httptrace.ClientTrace{
		GetConn: func(hostPort string) {
			config.logger.Infof("logId:%v,Prepare to get a connection for %s.", config.logId, hostPort)
		},
		GotConn: func(info httptrace.GotConnInfo) {
			config.logger.Infof("logId:%v,Got a connection: reused: %v, from the idle pool: %v.",
				config.logId, info.Reused, info.WasIdle)
		},
		PutIdleConn: func(err error) {
			if err == nil {
				config.logger.Infof("logId:%v,Put a connection to the idle pool: ok.", config.logId)
			} else {
				config.logger.Infof("logId:%v,Put a connection to the idle pool:%v", config.logId, err.Error())
			}
		},
		DNSStart: func(dnsStartInfo httptrace.DNSStartInfo) {
			config.logger.Infof("logId:%v,Begin DNS lookup for %v.", config.logId, dnsStartInfo.Host)
		},
		DNSDone: func(dnsDoneInfo httptrace.DNSDoneInfo) {
			config.logger.Infof("logId:%v,End DNS lookup, detail:%+v\n", config.logId, dnsDoneInfo)
		},
		ConnectStart: func(network, addr string) {
			config.logger.Infof("logId:%v,Dialing... (%s:%s).", config.logId, network, addr)
		},
		ConnectDone: func(network, addr string, err error) {
			if err == nil {
				config.logger.Infof("logId:%v,Dial is done. (%s:%s)", config.logId, network, addr)
			} else {
				config.logger.Infof("logId:%v,Dial is done with error: %s. (%s:%s)", config.logId, err, network, addr)
			}
		},
		TLSHandshakeStart: func() {
			config.logger.Infof("logId:%v,Begin TLSHandshake.", config.logId)
		},
		TLSHandshakeDone: func(connectionState tls.ConnectionState, i error) {
			config.logger.Infof("logId:%v,End TLSHandshake, detail:%+v", config.logId, connectionState)
		},
		WroteHeaders: func() {
			config.logger.Infof("logId:%v,Wrote headers: ok.", config.logId)
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			if info.Err == nil {
				config.logger.Infof("logId:%v,Wrote a request: ok.", config.logId)
			} else {
				config.logger.Infof("logId:%v,Wrote a request:%v", config.logId, info.Err.Error())
			}
		},
		GotFirstResponseByte: func() {
			config.logger.Infof("logId:%v,Got the first response byte.", config.logId)
		},
	})
	return traceCtx
}

func genReq(method Method, target string, config *Config) (*http.Request, error) {
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

	req, err := http.NewRequest(
		method.String(),
		urlEncode,
		bytes.NewReader(config.body),
	)
	if err != nil {
		return nil, err
	}

	for k, v := range config.header {
		req.Header.Set(k, v)
	}

	config.logger.Infow("log req",
		"id", config.logId,
		"method", req.Method,
		"urlEncode", req.RequestURI,
		"body", func() string {
			if len(config.body) > defaultBodyVerbose {
				return fmt.Sprintf("%v...", string(config.body[:defaultBodyVerbose]))
			}
			return string(config.body)
		}(),
		"header", fmt.Sprintf("%v", req.Header),
	)

	if config.trace {
		req = req.WithContext(genTraceCtx(config))
	} else {
		req = req.WithContext(config.context)
	}

	return req, nil
}

func Do(ctx context.Context, method Method, target string, opts ...Option) ([]byte, error) {
	config := newDefaultHttpConfig()

	{
		for _, opt := range opts {
			opt.apply(config)
		}
		config.context = ctx
	}

	if err := config.context.Err(); err != nil {
		config.logger.Errorw("context.Err not nil", "id", config.logId, "err", err)
		return nil, err
	}

	do := func(req *http.Request) ([]byte, error) {
		resp, err := config.client.Do(req)
		if err != nil {
			config.logger.Errorw("client.do failed", "id", config.logId, "err", err)
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}

	var rst []byte
	err := retry.Do(
		func() error {
			req, reqErr := genReq(method, target, config)
			if reqErr != nil {
				config.logger.Errorf("logId:%v,http.Do->genReq failed,err:%v", config.logId, reqErr)
				return reqErr
			}
			body, err := do(req)
			if err != nil {
				config.logger.Errorf("logId:%v,http.Do->do failed,err:%v", config.logId, err)
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
