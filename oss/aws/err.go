package aws

import "fmt"

/*
	go get golang.org/x/tools/cmd/stringer
*/

//go:generate stringer -type=Err  -linecomment
type err int

func errWrapper(err err) error {
	return fmt.Errorf("code:%d,errDesc:%s", err, err.String())
}

const (
	UnknownErrCode err = iota
	IllegalParams      // wrong params
)
