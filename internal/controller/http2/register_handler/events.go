package register_handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

// Events регистрирует обработчик для получения потока событий.
//
// Метод: GET /events
func Events(router *fiber.App, uc UsecasesForEvents, eventListener eventListener) {
	router.Get(
		"/events",
		func(ctx *fiber.Ctx) error {
			session := Session(ctx)
			ctx2, cancel := context.WithTimeout(ctx.Context(), time.Millisecond*100)
			defer cancel()
			return eventListener.Listen(ctx2, session.UserID, session.ID, func(event any) {
				ctx.Set("Content-Type", "text/event-stream")
				ctx.Set("Cache-Control", "no-cache")
				ctx.Set("Connection", "keep-alive")
				ctx.JSON(event)
			})
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(uc),
	)
}

// UsecasesForAcceptInvitation определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForEvents interface {
	middleware.UsecasesForRequireAuthorizedSession
}

type eventListener interface {
	Listen(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, f func(event any)) error
}
