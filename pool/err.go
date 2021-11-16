package pool

import "fmt"

//go:generate stringer -type=Err  -linecomment
type Err int

func ErrWrapper(err Err, ps ...string) error {
	return fmt.Errorf("code:%d,errDesc:%s,ps:%v", err, err.String(), ps)
}

const (
	UnknownErrCode    Err = iota
	IllegalParams         // illegal params
	GetConnTimeout        // Get Connection Timeout
	ResourceExhausted     // Pool Was Exhausted
)
