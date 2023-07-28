//go:build wireinject

//go:generate wire

package basics

import (
	"context"
	"github.com/google/wire"
)

func initializeBaz(ctx context.Context) (Baz, error) {
	panic(wire.Build(Set))
}
