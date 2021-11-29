package vars

import "fmt"

//go:generate stringer -type=Err  -linecomment
type Err int

func ErrWrapper(err Err, ps ...string) error {
	return fmt.Errorf("code:%d,errDesc:%s,ps:%v", err, err.String(), ps)
}

const (
	UnknownErrCode       Err = iota
	PluginMethodNotFound     // NewPlugin not found in plug
	ServerMethodNotFound     // NewServer not found in acceptor
)
