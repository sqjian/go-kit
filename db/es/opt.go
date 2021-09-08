package es

import "net/http"

type Option interface {
	apply(*Cli)
}

type optionFunc func(*Cli)

func (f optionFunc) apply(validator *Cli) {
	f(validator)
}

func WithHosts(hosts ...string) Option {
	return optionFunc(func(cli *Cli) {
		cli.meta.hosts = hosts
	})
}


func WithHttpClient(httpCli *http.Client) Option {
	return optionFunc(func(cli *Cli) {
		cli.meta.cli = httpCli
	})
}

func WithDebugInfo(debug bool) Option {
	return optionFunc(func(cli *Cli) {
		cli.meta.debug = debug
	})
}
