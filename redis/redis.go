package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var inst *redis.ClusterClient

func Init(addr []string, passwd string) error {
	if len(addr) == 0 {
		return fmt.Errorf("addr:%v is illegal", addr)
	}

	switch passwd {
	case "":
		{
			inst = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs: addr,
			})
		}
	default:
		{
			inst = redis.NewClusterClient(&redis.ClusterOptions{
				Addrs:    addr,
				Password: passwd,
			})
		}
	}
	if pingErr := inst.Ping(context.Background()).Err(); pingErr != nil {
		return pingErr
	}
	return nil
}

func HGet(ctx context.Context, key, field string) (string, error) {
	return inst.HGet(ctx, key, field).Result()
}

func HSet(ctx context.Context, key string, field, val string) error {
	return inst.HSet(ctx, key, field, val).Err()
}
