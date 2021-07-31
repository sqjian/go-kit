package uid

type Option interface {
	apply(*generator)
}

type optionFunc func(*generator)

func (f optionFunc) apply(log *generator) {
	f(log)
}

func WithNodeId(NodeId int64) Option {
	return optionFunc(func(generator *generator) {
		generator.MetaData.nodeId = NodeId
	})
}
