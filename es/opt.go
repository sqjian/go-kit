package es

import "net/http"

type OptionFunc func(*Cli)

func WithHosts(hosts ...string) OptionFunc {
	return func(cli *Cli) {
		cli.config.hosts = hosts
	}
}

func WithHttpClient(httpCli *http.Client) OptionFunc {
	return func(cli *Cli) {
		cli.config.cli = httpCli
	}
}

func WithDebugInfo(debug bool) OptionFunc {
	return func(cli *Cli) {
		cli.config.debug = debug
	}
}
