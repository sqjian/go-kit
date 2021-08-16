package preset

//go:generate stringer -type=LogType  -linecomment
type LogType int

const (
	_ LogType = iota
	Zap
)
