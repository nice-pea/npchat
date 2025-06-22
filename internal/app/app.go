package app

import (
	"context"

	"golang.org/x/sync/errgroup"
)

func Run(ctx context.Context, cfg Config) error {
	g, ctx := errgroup.WithContext(ctx)

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
