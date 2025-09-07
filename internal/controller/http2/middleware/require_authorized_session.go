package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	findSession "github.com/nice-pea/npchat/internal/usecases/sessions/find_session"
)

const CtxKeyUserSession = "userSession"

// RequireAuthorizedSession требует авторизованную сессии
func RequireAuthorizedSession(uc UsecasesForRequireAuthorizedSession, tm UsecasesForRequireAuthorizedParseJWT) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		// Прочитать заголовок
		header := ctx.Get("Authorization")

		parts := strings.Split(header, " ")
		if len(parts) != 2 {
			return fiber.ErrUnauthorized
		}

		authType := parts[0]
		token := parts[1]

		switch authType {
		case "Bearer":
			out, err := bearer(uc, token)
			if err != nil {
				return err
			}
			// Сохранить сессию в контекст
			ctx.Locals(CtxKeyUserSession, out.Sessions[0])

		case "JWT":
			out, err := parseJwt(tm, token)
			if err != nil {
				return err
			}
			ctx.Locals("UserID", out.UserID)
			ctx.Locals("SessionID", out.SessionID)

		default:
			return fiber.ErrUnauthorized
		}

		return ctx.Next()
	}
}

// UsecasesForRequireAuthorizedSession определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForRequireAuthorizedSession interface {
	FindSessions(findSession.In) (findSession.Out, error)
}

func bearer(uc UsecasesForRequireAuthorizedSession, token string) (findSession.Out, error) {
	// Найти сессию по токену
	out, err := uc.FindSessions(findSession.In{
		Token: token,
	})
	if err != nil {
		return findSession.Out{}, fiber.ErrInternalServerError
	}
	if len(out.Sessions) != 1 {
		return findSession.Out{}, fiber.ErrUnauthorized
	}

	return out, nil
}

type OutJWT struct {
	UserID    string
	SessionID string
}

type UsecasesForRequireAuthorizedParseJWT interface {
	Parse(token string) (OutJWT, error)
}

func parseJwt(ucjwt UsecasesForRequireAuthorizedParseJWT, token string) (OutJWT, error) {
	out, err := ucjwt.Parse(token)
	if err != nil {
		return OutJWT{}, fiber.ErrUnauthorized
	}

	return out, nil
}
