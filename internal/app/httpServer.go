package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/registerHandler"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/router"
)

func initHttpServer(ss *services) *http.Server {
	r := &router.Router{
		Services: ss,
	}
	registerHandlers(r)

	return &http.Server{
		Addr:    ":8080",
		Handler: r,
	}
}

func runHttpServer(ctx context.Context, server *http.Server) error {
	g, ctx := errgroup.WithContext(ctx)

	// Запуск сервера
	g.Go(func() error {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server.ListenAndServe: %w", err)
		}
		return nil
	})

	// Завершение сервера при завершении контекста
	g.Go(func() error {
		// > The first call to return a non-nil error cancels the group's context
		<-ctx.Done()
		return server.Shutdown(ctx)
	})

	return g.Wait()
}

func registerHandlers(r http2.Router) {
	// Служебные
	registerHandler.Ping(r)

	// Аутентификация
	registerHandler.LoginByPassword(r)

	// Чат
	registerHandler.MyChats(r)
	registerHandler.CreateChat(r)
	registerHandler.UpdateChatName(r)

	// Участники
	registerHandler.LeaveChat(r)
	registerHandler.ChatMembers(r)
	registerHandler.DeleteMember(r)

	// Приглашение
	registerHandler.MyInvitations(r)
	registerHandler.ChatInvitations(r)

	// Управление приглашениями
	registerHandler.SendInvitation(r)
	registerHandler.AcceptInvitation(r)
	registerHandler.CancelInvitation(r)
}
