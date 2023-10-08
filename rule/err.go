package rule

import "fmt"

/*
	go get golang.org/x/tools/cmd/stringer
*/

//go:generate stringer -type=Err  -linecomment
type Err int

func errWrapper(err Err) error {
	return fmt.Errorf("code:%d,errDesc:%s", err, err.String())
}

const (
	UnknownErrCode Err = iota
	NotFound
)
