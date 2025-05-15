package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/sync/errgroup"

	"github.com/saime-0/nice-pea-chat/internal/controller/handler"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/router"
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
	handler.RegisterPingHandler(r)

	// Аутентификация
	handler.RegisterLoginByPasswordHandler(r)

	// Чат
	handler.RegisterMyChatsHandler(r)
	handler.RegisterCreateChatHandler(r)
	handler.RegisterUpdateChatNameHandler(r)

	// Участники
	handler.RegisterLeaveChatHandler(r)
	handler.RegisterChatMembersHandler(r)
	handler.RegisterDeleteMemberHandler(r)

	// Приглашение
	handler.RegisterMyInvitationsHandler(r)
	handler.RegisterChatInvitationsHandler(r)

	// Управление приглашениями
	handler.RegisterSendInvitationHandler(r)
	handler.RegisterAcceptInvitationHandler(r)
	handler.RegisterCancelInvitationHandler(r)
}
