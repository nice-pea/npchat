package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/saime-0/nice-pea-chat/internal/controller/handler"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
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

	r := &router.Router{
		Services: ss,
		Middlewares: []http2.Middleware{
			middleware.RequireRequestID,
			middleware.RequireAcceptJson,
			middleware.RequireAuthorizedSession,
		},
	}
	handler.Ping(r)
	handler.CreateChat(r)
	rau := &router.Router{
		Services: ss,
	}
	handler.LoginByPassword(rau)

	m := http.ServeMux{}
	m.

		// Инициализация контроллера
		//r := controller.InitController(ss.chats, ss.invitations, ss.members, ss.sessions, ss.authnPassword)

		// Запуск сервера
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
