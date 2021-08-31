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

type Plugin interface {
	Name() string
	Init() error
	FInit() error
	Interest(msg *proto.Msg, tools PluginTools) (Action, error)
	PreProcess(msg *proto.Msg, tools PluginTools) (Action, error)
	Process(msg *proto.Msg, tools PluginTools) (Action, error)
	PostProcess(msg *proto.Msg, tools PluginTools) (Action, error)
}
