//go:build wireinject

//go:generate wire

package _struct

import "github.com/google/wire"

func injectFooBar() FooBar {
	panic(wire.Build(Set))
}
