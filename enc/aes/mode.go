package aes

//go:generate stringer -type=Mode  -linecomment
type Mode int64

const (
	ECB Mode = iota
	CBC
	CFB
	CTR
	OFB
)
