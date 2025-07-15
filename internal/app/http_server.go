package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/sync/errgroup"

	registerHandler "github.com/nice-pea/npchat/internal/controller/http2/register_handler"
)

func runHttpServer(ctx context.Context, ss *services, cfg Config) error {
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: fiberErrorHandler,
	})
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

// registerHandlers регистрирует обработчики
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

// fiberErrorHandler разделяет составные ошибки и помещает в body.
func fiberErrorHandler(ctx *fiber.Ctx, err error) error {
	var resp errorResponse
	for _, err2 := range errorsFlatten(err) {
		resp.Errors = append(resp.Errors, errAsMap(err2))
	}
	return ctx.JSON(resp)
}

// errAsMap преобразует ошибку в карту
func errAsMap(err error) map[string]any {
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		return map[string]any{
			"code":    strconv.Itoa(fiberErr.Code),
			"message": fiberErr.Message,
		}
	}

	return map[string]any{
		//"code": "", // Пока не продумал полностью детализацию и коды для ошибок
		"message": err.Error(),
	}
}

// errorsFlatten рекурсивно извлекает все ошибки из составной ошибки.
func errorsFlatten(err error) []error {
	var errs []error
	if err == nil {
		return errs
	}

	// Проверяем, поддерживает ли ошибка интерфейс Unwrap() []error
	if unwrap, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range unwrap.Unwrap() {
			errs = append(errs, errorsFlatten(e)...)
		}
	} else {
		// Простая ошибка — добавляем её
		errs = append(errs, err)
	}

	return errs
}

type errorResponse struct {
	Errors []map[string]any `json:"errors"`
}
