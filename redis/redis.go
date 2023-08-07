package redis

import (
	"context"
	"time"
)

type Redis interface {
	Set(ctx context.Context, key string, value any, duration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	HSet(ctx context.Context, name string, values map[string]any) error
	HGet(ctx context.Context, name, field string) (string, error)
	HGetAll(ctx context.Context, name string) (map[string]string, error)
}
