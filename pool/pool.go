package pool

import (
	"context"
)

func NewPool(ctx context.Context, opts ...OptionFunc) (Pool, error) {
	cfg := newDefaultCfg()
	for _, opt := range opts {
		opt(cfg)
	}

	switch cfg.PoolType {
	case Exclusive:
		return NewExclusivePool(ctx, cfg)
	case Share:
		return NewSharePool(ctx, cfg)
	default:
		return nil, ErrWrapper(IllegalParams)
	}
}
