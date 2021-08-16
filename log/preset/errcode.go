package preset

import "fmt"

/*
	go get golang.org/x/tools/cmd/stringer
*/

//go:generate stringer -type=ErrCode  -linecomment
type ErrCode int

func ErrWrapper(errCode ErrCode) error {
	return fmt.Errorf("errCode:%d,errDesc:%s", errCode, errCode.String())
}

const (
	UnknownErrCode ErrCode = iota
	IllegalParams
	IllegalKeyType
)
