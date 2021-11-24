package pool

//go:generate stringer -type=Err  -linecomment
type Type int

const (
	Exclusive Type = iota
	Share
)
