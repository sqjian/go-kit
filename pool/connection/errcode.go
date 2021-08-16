package connection

import "fmt"

//go:generate stringer -type=ErrCode  -linecomment
type ErrCode int

func GenErr(errCode ErrCode) error {
	return fmt.Errorf("errCode:%d,errDesc:%s", errCode, errCode.String())
}

const (
	UnknownErrCode ErrCode = iota
	IllegalParams          // illegal params
	GetConnTimeout         // Get Connection Timeout
	PoolExhausted          // Pool Was Exhausted
)
