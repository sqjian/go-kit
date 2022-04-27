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

func newDefaultCliCfg() *cliCfg {
	return &cliCfg{
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

func genTraceCtx(config *cliCfg) context.Context {
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

func genReq(method Method, target string, cfg *cliCfg) (*http.Request, error) {
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

	{
		cfg.logger.Infof("log req => id:%v,method:%v,urlEncode:%v", cfg.logId, req.Method, req.RemoteAddr)
		cfg.logger.Infof("log req => id:%v,header:%v", cfg.logId, req.Header)
		cfg.logger.Infof("log req => id:%v,body:%v", cfg.logId, func() string {
			if len(cfg.body) > defaultBodyVerbose {
				return fmt.Sprintf("%v...", string(cfg.body[:defaultBodyVerbose]))
			}
			return string(cfg.body)
		}())
	}

	if cfg.trace {
		req = req.WithContext(genTraceCtx(cfg))
	} else {
		req = req.WithContext(cfg.context)
	}

	return req, nil
}

func Do(ctx context.Context, method Method, target string, opts ...CliOption) ([]byte, error) {
	cfg := newDefaultCliCfg()
	{
		for _, opt := range opts {
			opt.apply(cfg)
		}
		cfg.context = ctx
	}

	if err := cfg.context.Err(); err != nil {
		cfg.logger.Errorf("context.Err not nil => id:%v,err:%v", cfg.logId, err)
		return nil, err
	}

	do := func(req *http.Request) ([]byte, error) {
		resp, err := cfg.client.Do(req)
		if err != nil {
			cfg.logger.Errorf("client.do failed => id:%v,err:%v", cfg.logId, err)
			return nil, err
		}
		defer resp.Body.Close()
		respData, respDataErr := ioutil.ReadAll(resp.Body)
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
				cfg.logger.Errorf("logId:%v,http.Do->genReq failed,err:%v", cfg.logId, reqErr)
				return reqErr
			}
			body, err := do(req)
			if err != nil {
				cfg.logger.Errorf("logId:%v,http.Do->do failed,err:%v", cfg.logId, err)
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
