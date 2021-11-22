package mq

import "fmt"

/*
	go get golang.org/x/tools/cmd/stringer
*/

//go:generate stringer -type=Err  -linecomment
type Err int

func ErrWrapper(err Err) error {
	return fmt.Errorf("code:%d,errDesc:%s", err, err.String())
}

const (
	ErrTopicFormat    Err = iota // pulsar should in format (persistent|non-persistent)://tenant/namespace/topic
	ErrInvalidUrl                // invalid mq url
	ErrNoSuchConsumer            // no such consumer
	ErrConnectBroken             // mq net connection is broken
	ErrDeadline                  // context deadline exceed
	ErrMqType                    // invalid mq type
)
