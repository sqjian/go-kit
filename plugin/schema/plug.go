package schema

import (
	"github.com/sqjian/go-kit/plugin/proto"
)

type Action int

const (
	STOP Action = iota
	CONTINUE
	SKIP
)

type Plug interface {
	Name() string
	Init() error
	FInit() error
	Interest(*proto.Msg, PlugTools) (Action, error)
	PreProcess(*proto.Msg, PlugTools) (Action, error)
	Process(*proto.Msg, PlugTools) (Action, error)
	PostProcess(*proto.Msg, PlugTools) (Action, error)
}
