package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/controller/handler"
	"github.com/saime-0/nice-pea-chat/internal/controller/router"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

func Run(ctx context.Context) error {
	// Инициализация репозиториев
	repos, closer, err := initSqliteRepositories(sqlite.Config{
		MigrationsDir: "migrations/repository/sqlite",
	})
	if err != nil {
		return err
	}
	defer closer()

	// Инициализация сервисов
	ss := initServices(repos)

	// Инициализация http контроллера
	r := &router.Router{
		Services: ss,
	}
	handler.RegisterPingHandler(r)
	handler.CreateChat(r)
	handler.LoginByPassword(r)

	// Запуск http сервера
	server := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
	errChan := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("server.ListenAndServe: %w", err)
		}
	}()

	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("server.Shutdown: %w", err)
		}
		return nil
	}
}
