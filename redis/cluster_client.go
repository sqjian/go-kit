package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func newClusterClient(addrs []string, passwd string) (Redis, error) {
	if len(addrs) == 0 {
		return nil, fmt.Errorf("addrs:%v is illegal", addrs)
	}

	inst := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: passwd,
	})

	if pingErr := inst.Ping(context.Background()).Err(); pingErr != nil {
		return nil, pingErr
	}

	return &client{ClusterClient: inst}, nil
}

type client struct {
	*redis.ClusterClient
}

func (c *client) Set(ctx context.Context, key string, value any, duration time.Duration) error {
	return c.ClusterClient.Set(ctx, key, value, duration).Err()
}

func (c *client) Get(ctx context.Context, key string) (string, error) {
	return c.ClusterClient.Get(ctx, key).Result()
}

func (c *client) HSet(ctx context.Context, name string, values map[string]any) error {
	return c.ClusterClient.HSet(ctx, name, values).Err()
}

func (c *client) HGet(ctx context.Context, name, field string) (string, error) {
	return c.ClusterClient.HGet(ctx, name, field).Result()
}

func (c *client) HGetAll(ctx context.Context, name string) (map[string]string, error) {
	return c.ClusterClient.HGetAll(ctx, name).Result()
}
