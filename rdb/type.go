package rdb

//go:generate stringer -type=Type  -linecomment
type Type int64

const (
	UnknownType Type = iota
	Mysql
	Sqlite
)
