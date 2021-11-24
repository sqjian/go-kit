package log

//go:generate stringer -type=Level  -linecomment
type Level int64

const (
	UnknownLevel Level = iota
	Dummy
	Debug
	Info
	Warn
	Error
)
