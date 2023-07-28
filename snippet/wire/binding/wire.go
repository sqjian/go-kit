//go:build wireinject

//go:generate wire

package binding

import "github.com/google/wire"

func injectFooBar() Footer {
	panic(wire.Build(Set))
}
