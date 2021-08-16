package connection

import "fmt"

//go:generate stringer -type=Err  -linecomment
type Err int

func ErrWrapper(err Err) error {
	return fmt.Errorf("err:%d,errDesc:%s", err, err.String())
}

const (
	UnknownErrCode Err = iota
	IllegalParams      // illegal params
	GetConnTimeout     // Get Connection Timeout
	PoolExhausted      // Pool Was Exhausted
)
