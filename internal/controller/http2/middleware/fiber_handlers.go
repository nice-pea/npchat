package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/nice-pea/npchat/internal/service"
)

func RequareAuthoruzation(sessions *service.Sessions) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Прочитать заголовок
		header := ctx.Get("Authorization")
		token, _ := strings.CutPrefix(header, "Bearer ")
		if token == "" {
			return fiber.ErrUnauthorized
		}

		// Найти сессию по токену
		userSessions, err := sessions.Find(service.SessionsFindIn{
			Token: token,
		})
		if err != nil {
			return fiber.ErrInternalServerError
		}
		if len(userSessions) != 1 {
			return fiber.ErrUnauthorized
		}

		// Сохранить сессию в контекст
		ctx.Locals(CtxUserSession, userSessions[0])

		return nil
	}
}

var CtxUserSession = "userSession"
