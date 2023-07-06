package es

import "net/http"

type OptionFunc func(*cli)

func WithHosts(hosts ...string) OptionFunc {
	return func(cli *cli) {
		cli.config.hosts = hosts
	}
}

func WithHttpClient(httpCli *http.Client) OptionFunc {
	return func(cli *cli) {
		cli.config.cli = httpCli
	}
}

func WithDebugInfo(debug bool) OptionFunc {
	return func(cli *cli) {
		cli.config.debug = debug
	}
}
