package err

import "fmt"

/*
	go install golang.org/x/tools/cmd/stringer
*/

//go:generate stringer -type=Err  -linecomment
type Err int

func Wrapper(err Err) error {
	return fmt.Errorf("code:%d,errDesc:%s", err, err.String())
}

const (
	UnknownCode Err = iota
	Code1           //this is error code1
	Code2           //this is error code2
	Code3           //this is error code3
)
