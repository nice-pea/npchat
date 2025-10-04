package redisRegistry

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Config конфигурация для подключения к Redis
type Config struct {
	DSN string
}

// Init инициализирует и возвращает клиент Redis на основе указанной конфигурации
func Init(cfg Config) (*redis.Client, error) {
	// Парсинг строки подключения
	opt, err := redis.ParseURL(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("redis.ParseURL: %w", err)
	}

	// Создание клиента Redis
	client := redis.NewClient(opt)

	// Проверка соединения с Redis
	_, err = client.Ping(context.Background()).Result()
	if err != nil {
		return nil, fmt.Errorf("redis.Ping: %w", err)
	}

	return client, nil
}
