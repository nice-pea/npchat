package app

import (
	"context"

	"golang.org/x/sync/errgroup"

	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

func Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	// Инициализация репозиториев
	repos, closeSqliteRepos, err := initSqliteRepositories(sqlite.Config{
		MigrationsDir: "migrations/repository/sqlite",
	})
	if err != nil {
		return err
	}
	defer closeSqliteRepos()

	// Инициализация сервисов
	ss := initServices(repos)

	// Инициализация и Запуск http контроллера
	server := initHttpServer(ss)
	g.Go(func() error {
		return runHttpServer(ctx, server)
	})

	return g.Wait()
}
