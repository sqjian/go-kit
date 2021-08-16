package errcode

import "fmt"

/*
	go get golang.org/x/tools/cmd/stringer
*/

//go:generate stringer -type=ErrCode  -linecomment
type ErrCode int

func GenErr(errCode ErrCode) error {
	return fmt.Errorf("errCode:%d,errDesc:%s", errCode, errCode.String())
}

const (
	UnknownCode ErrCode = iota
	Code1               //this is error code1
	Code2               //this is error code2
	Code3               //this is error code3
)
