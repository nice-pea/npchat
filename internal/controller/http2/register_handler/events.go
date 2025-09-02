package register_handler

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

// Events регистрирует обработчик для получения потока событий.
//
// Метод: GET /events
func Events(router *fiber.App, uc UsecasesForEvents, eventListener eventListener) {
	router.Get(
		"/events",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc),
		func(ctx *fiber.Ctx) error {
			ctx.Set("Content-Type", "text/event-stream")
			ctx.Set("Cache-Control", "no-cache")
			ctx.Set("Connection", "keep-alive")

			session := Session(ctx)

			// Контекст для корректной остановки прослушивания событий
			ctx2, cancel := context.WithCancel(context.Background())
			// Канал для отслеживания завершения запроса
			done := ctx.Context().Done()
			// Отменить контекст при завершении запроса
			go func() {
				<-done
				cancel()
			}()

			// Регистрация обработчика для отправки потока сообщений
			ctx.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
				err := eventListener.Listen(ctx2, session.UserID, session.ID, func(event any) {
					// Отправить данные в буфер
					if err := json.NewEncoder(w).Encode(event); err != nil {
						slog.Warn("eventListener.Listen: json.NewEncoder: " + err.Error())
					}
					// Отправить данные во writer и очистить буфер
					if err := w.Flush(); err != nil {
						slog.Warn("eventListener.Listen: w.Flush: " + err.Error())
					}
				})
				if err != nil && !errors.Is(err, context.Canceled) {
					slog.Warn("eventListener.Listen:" + err.Error())
				}
			}))

			return nil
		},
	)
}

// UsecasesForAcceptInvitation определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForEvents interface {
	middleware.UsecasesForRequireAuthorizedSession
}

type eventListener interface {
	Listen(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, f func(event any)) error
}
