package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	sessionsFind "github.com/nice-pea/npchat/internal/service/sessions_find"
)

// RequireAuthorizedSession требует авторизованную сессии
func RequireAuthorizedSession(uc interface {
	SessionsFind(sessionsFind.In) (sessionsFind.Out, error)
}) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Прочитать заголовок
		header := ctx.Get("Authorization")
		token, _ := strings.CutPrefix(header, "Bearer ")
		if token == "" {
			return fiber.ErrUnauthorized
		}

		// Найти сессию по токену
		out, err := uc.SessionsFind(sessionsFind.In{
			Token: token,
		})
		if err != nil {
			return fiber.ErrInternalServerError
		}
		if len(out.Sessions) != 1 {
			return fiber.ErrUnauthorized
		}

		// Сохранить сессию в контекст
		ctx.Locals(CtxUserSession, out.Sessions[0])

		return nil
	}
}

var CtxUserSession = "userSession"
