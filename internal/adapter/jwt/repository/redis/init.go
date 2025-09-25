package redisCache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	DSN string
}

func Init(cfg Config) (*redis.Client, error) {
	opt, err := redis.ParseURL(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("redisCache.Init: %w", err)
	}

	client := redis.NewClient(opt)

	ctx := context.Background()
	_, err = client.Ping(ctx).Result()

	if err != nil {
		return nil, fmt.Errorf("redisCache.Init: %w", err)
	}

	return client, nil
}
