package easyhttp

//go:generate stringer -type=Method  -linecomment
type Method int64

const (
	GET Method = iota
	POST
)
