package digest

import "fmt"

/*
	go get golang.org/x/tools/cmd/stringer
*/

//go:generate stringer -type=Err  -linecomment
type Err int

func ErrWrapper(err Err) error {
	return fmt.Errorf("err:%d,errDesc:%s", err, err.String())
}

const (
	UnknownErrCode Err = iota
	IllegalKeyType     //Illegal KeyType
)
