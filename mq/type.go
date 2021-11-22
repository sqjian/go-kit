package mq

//go:generate stringer -type=Type  -linecomment
type Type int

const (
	_ Type = iota
	Pulsar
)
