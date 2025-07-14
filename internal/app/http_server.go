package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"

	registerHandler "github.com/nice-pea/npchat/internal/controller/http2/register_handler"
)

func runHttpServer(ctx context.Context, ss *services, cfg Config) error {
	fiberApp := fiber.New()
	registerHandlers(fiberApp, ss)

	g, ctx := errgroup.WithContext(ctx)

	// Запуск сервера
	g.Go(func() error {
		err := fiberApp.Listen(cfg.HttpAddr)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("server.ListenAndServe: %w", err)
		}
		return nil
	})

	// Завершение сервера при завершении контекста
	g.Go(func() error {
		<-ctx.Done()
		return fiberApp.Shutdown()
	})

	return g.Wait()
}

func registerHandlers(r *fiber.App, ss *services) {
	// Служебные
	registerHandler.Ping(r, ss)

	// OAuth /oauth
	registerHandler.OAuthInitRegistration(r, ss)
	registerHandler.OAuthCompleteRegistrationCallback(r, ss)

	// Аутентификация /auth
	registerHandler.LoginByPassword(r, ss)
	registerHandler.RegistrationByPassword(r, ss)

	// Чат /chats
	registerHandler.MyChats(r, ss)
	registerHandler.CreateChat(r, ss)
	registerHandler.UpdateChatName(r, ss)
	registerHandler.LeaveChat(r, ss)
	registerHandler.ChatMembers(r, ss)
	registerHandler.ChatInvitations(r, ss)

	// Участники /chats//members
	registerHandler.DeleteMember(r, ss)

	// Приглашения /invitations
	registerHandler.MyInvitations(r, ss)
	registerHandler.SendInvitation(r, ss)
	registerHandler.AcceptInvitation(r, ss)
	registerHandler.CancelInvitation(r, ss)
}
