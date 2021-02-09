package errors

import "fmt"

var (
	InternalErr = New(-1, "internal error")
)

type Error interface {
	Error() string
	ErrorCode() int
	ErrorMsg() string
}

type ErrorPredefined struct {
	Msg  string
	Code int
}

func New(code int, msg string) *ErrorPredefined {
	return &ErrorPredefined{
		Code: code,
		Msg:  msg,
	}
}

func (e *ErrorPredefined) Error() string {
	return fmt.Sprintf("error: %s | Code: %d", e.Msg, e.Code)
}
func (e *ErrorPredefined) ErrorCode() int {
	return e.Code
}
func (e *ErrorPredefined) ErrorMsg() string {
	return e.Msg
}
