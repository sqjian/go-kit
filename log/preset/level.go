package preset

//go:generate stringer -type=Level  -linecomment
type Level int64

const (
	UnknownLevel Level = iota
	None
	Debug
	Info
	Warn
	Error
)
