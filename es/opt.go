package es

import "net/http"

type OptionFunc func(*Cli)

func WithHosts(hosts ...string) OptionFunc {
	return func(cli *Cli) {
		cli.meta.hosts = hosts
	}
}

func WithHttpClient(httpCli *http.Client) OptionFunc {
	return func(cli *Cli) {
		cli.meta.cli = httpCli
	}
}

func WithDebugInfo(debug bool) OptionFunc {
	return func(cli *Cli) {
		cli.meta.debug = debug
	}
}
