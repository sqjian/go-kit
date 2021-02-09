package http

import (
	"bytes"
	"context"
	"github.com/sqjian/go-kit/retry"
	"io/ioutil"
	"net"
	"net/http"
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
	defaultMaxConnsPerHost     = 10
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
			MaxConnsPerHost:     defaultMaxConnsPerHost,
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
	}
}
func Do(
	target string,
	method Method,
	query map[string]string,
	header map[string]string,
	data []byte,
	opts ...Option,
) ([]byte, error) {
	config := newDefaultHttpConfig()

	for _, opt := range opts {
		opt.apply(config)
	}

	if err := config.context.Err(); err != nil {
		return nil, err
	}

	httpSend := func() ([]byte, error) {
		u, err := url.Parse(target)
		if err != nil {
			return nil, err
		}

		q := u.Query()
		for k, v := range query {
			q.Set(k, v)
		}
		u.RawQuery = q.Encode()
		urlEncode := u.String()

		req, err := http.NewRequest(method.String(), urlEncode, bytes.NewReader(data))
		if err != nil {
			return nil, err
		}
		req = req.WithContext(config.context)

		resp, err := config.client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		return ioutil.ReadAll(resp.Body)
	}

	var rst []byte
	err := retry.Do(
		func() error {
			body, err := httpSend()
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
