package pool

import (
	"context"
)

func NewPool(ctx context.Context, opts ...Option) (Pool, error) {
	cfg := newDefaultCfg()
	for _, opt := range opts {
		opt.apply(cfg)
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
