package pulsar

type Option interface {
	apply(*Client)
}

type optionFunc func(*Client)

func (f optionFunc) apply(cli *Client) {
	f(cli)
}

func WithDebug(debug bool) Option {
	return optionFunc(func(cli *Client) {
		cli.meta.debug = debug
	})
}
