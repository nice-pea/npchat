package app

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg Config) error {
	g, ctx := errgroup.WithContext(ctx)

	slog.SetLogLoggerLevel(cfg.SlogLevel)
	slog.Info(fmt.Sprintf("Уровень логирования: %s", cfg.SlogLevel))

	// Инициализация репозиториев
	repos, closeSqliteRepos, err := initSqliteRepositories(cfg.SQLite)
	if err != nil {
		return err
	}
	defer closeSqliteRepos()

	// Инициализация адаптеров
	adaps := initAdapters(cfg)

	// Инициализация сервисов
	ss := initServices(repos, adaps)

	// Инициализация и Запуск http контроллера
	server := initHttpServer(ss, cfg)
	g.Go(func() error {
		return runHttpServer(ctx, server)
	})

	return g.Wait()
}
