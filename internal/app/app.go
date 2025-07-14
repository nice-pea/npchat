package app

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg Config) error {
	g, ctx := errgroup.WithContext(ctx)

	slog.SetLogLoggerLevel(slogLevel(cfg.LogLevel))
	slog.Info(fmt.Sprintf("Уровень логирования: %s", cfg.LogLevel))

	// Инициализация репозиториев
	repos, closeRepos, err := initPgsqlRepositories(cfg.Pgsql)
	if err != nil {
		return err
	}
	defer closeRepos()

	// Инициализация адаптеров
	adaps := initAdapters(cfg)

	// Инициализация сервисов
	ss := initServices(repos, adaps)

	// Инициализация и Запуск http контроллера
	g.Go(func() error {
		return runHttpServer(ctx, ss, cfg)
	})

	return g.Wait()
}
