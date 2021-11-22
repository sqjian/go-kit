package mq

type Option interface {
	apply(*Instance)
}

type optionFunc func(*Instance)

func (f optionFunc) apply(inst *Instance) {
	f(inst)
}
