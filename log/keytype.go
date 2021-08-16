package log

//go:generate stringer -type=KeyType  -linecomment
type KeyType int

const (
	UnknownKeyType KeyType = iota
	Zap
)
