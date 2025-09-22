package redisCache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Addr     string
	Password string
	DB       int
}

func Init(cfg Config) (*redis.Client, error) {

	client := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	_, err := client.Ping(ctx).Result()

	if err != nil {
		return nil, fmt.Errorf("redisCache.Init: %w", err)
	}

	return client, nil
}
