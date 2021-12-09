package hash

//go:generate stringer -type=KeyType  -linecomment
type KeyType int

const (
	UnknownKeyType KeyType = iota
	MD5
	SHA1
)
