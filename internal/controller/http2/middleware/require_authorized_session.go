package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	findSession "github.com/nice-pea/npchat/internal/service/sessions/find_session"
)

const CtxKeyUserSession = "userSession"

// RequireAuthorizedSession требует авторизованную сессии
func RequireAuthorizedSession(uc UsecasesForRequireAuthorizedSession) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Прочитать заголовок
		header := ctx.Get("Authorization")
		token, _ := strings.CutPrefix(header, "Bearer ")
		if token == "" {
			return fiber.ErrUnauthorized
		}

		// Найти сессию по токену
		out, err := uc.FindSessions(findSession.In{
			Token: token,
		})
		if err != nil {
			return fiber.ErrInternalServerError
		}
		if len(out.Sessions) != 1 {
			return fiber.ErrUnauthorized
		}

		// Сохранить сессию в контекст
		ctx.Locals(CtxKeyUserSession, out.Sessions[0])

		return nil
	}
}

type UsecasesForRequireAuthorizedSession interface {
	FindSessions(findSession.In) (findSession.Out, error)
}
