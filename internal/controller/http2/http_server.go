package http2

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"golang.org/x/sync/errgroup"

	registerHandler "github.com/nice-pea/npchat/internal/controller/http2/register_handler"
)

type Config struct {
	HttpAddr string
}

// RunHttpServer запускает http сервер до момента отмена контекста
func RunHttpServer(ctx context.Context, uc RequiredUsecases, eventListener registerHandler.EventListener, requiredJwt RequiredJwt, cfg Config) error {
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: fiberErrorHandler,
	})
	registerHandlers(fiberApp, uc, requiredJwt, eventListener)

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
func registerHandlers(r *fiber.App, uc RequiredUsecases, jwtAuth RequiredJwt, eventListener registerHandler.EventListener) {
	// Подключение middleware для логирования
	r.Use(logger.New(logger.Config{
		TimeFormat: "2006-01-02 15:04:05",
	}))

	// Служебные
	registerHandler.Ping(r)

	registerHandler.Events(r, uc, eventListener, jwtAuth)

	// Oauth /oauth
	registerHandler.OauthAuthorize(r, uc)
	registerHandler.OauthCallback(r, uc, jwtAuth)

	// Аутентификация /auth
	registerHandler.LoginByPassword(r, uc, jwtAuth)
	registerHandler.RegistrationByPassword(r, uc, jwtAuth)

	// Чат /chats
	registerHandler.MyChats(r, uc, jwtAuth)
	registerHandler.CreateChat(r, uc, jwtAuth)
	registerHandler.UpdateChatName(r, uc, jwtAuth)
	registerHandler.LeaveChat(r, uc, jwtAuth)
	registerHandler.ChatMembers(r, uc, jwtAuth)
	registerHandler.ChatInvitations(r, uc, jwtAuth)

	// Участники /chats//members
	registerHandler.DeleteMember(r, uc, jwtAuth)

	// Приглашения /invitations
	registerHandler.MyInvitations(r, uc, jwtAuth)
	registerHandler.SendInvitation(r, uc, jwtAuth)
	registerHandler.AcceptInvitation(r, uc, jwtAuth)
	registerHandler.CancelInvitation(r, uc, jwtAuth)
}

// fiberErrorHandler разделяет составные ошибки и помещает в body.
func fiberErrorHandler(ctx *fiber.Ctx, err error) error {
	var errs []map[string]any
	for _, err2 := range errorsFlatten(err) {
		errs = append(errs, errAsMap(err2))
	}
	return ctx.JSON(fiber.Map{
		"errors": errs,
	})
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

	// https://habr.com/ru/companies/oleg-bunin/articles/913096/#:~:text=4.%20Wrapping%20%D0%B4%D0%BB%D1%8F%20%D0%BD%D0%B5%D1%81%D1%82%D1%80%D1%83%D0%BA%D1%82%D1%83%D1%80%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%BD%D1%8B%D1%85%20%D0%B4%D0%B0%D0%BD%D0%BD%D1%8B%D1%85
	var detailsErr interface{ Details() map[string]any }
	if errors.As(err, &detailsErr) {
		return map[string]any{
			"message": err.Error(),
			"details": detailsErr.Details(),
		}
	}

	return map[string]any{
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
