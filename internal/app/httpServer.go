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

	// Аутентификация /login
	registerHandler.LoginByPassword(r)

	// Чат /chats
	registerHandler.MyChats(r)
	registerHandler.CreateChat(r)
	registerHandler.UpdateChatName(r)
	registerHandler.LeaveChat(r)
	registerHandler.ChatMembers(r)
	registerHandler.ChatInvitations(r)

	// Участники /chats//members
	registerHandler.DeleteMember(r)

	// Приглашения /invitations
	registerHandler.MyInvitations(r)
	registerHandler.SendInvitation(r)
	registerHandler.AcceptInvitation(r)
	registerHandler.CancelInvitation(r)

	// OAuth /oauth
	registerHandler.GoogleRegistration(r)
	registerHandler.GoogleRegistrationCallback(r)
}
