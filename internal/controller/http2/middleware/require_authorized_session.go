package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"

	findSession "github.com/nice-pea/npchat/internal/usecases/sessions/find_session"
)

const (
	CtxKeyUserID    = "UserID"
	CtxKeySessionID = "SessionID"
)

const (
	SessionToken = "SessionToken"
	Bearer       = "Bearer"
)

// RequireAuthorizedSession требует авторизованную сессии
func RequireAuthorizedSession(uc UsecasesForRequireAuthorizedSession, jwtparser JwtParser) fiber.Handler {
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
		case SessionToken:
			out, err := findSessionf(uc, token)
			if err != nil {
				return err
			}
			// Сохранить сессию в контекст
			session := out.Sessions[0]
			ctx.Locals(CtxKeyUserID, session.UserID)
			ctx.Locals(CtxKeySessionID, session.ID)
		case Bearer:
			out, err := parseJwt(jwtparser, token)
			if err != nil {
				return err
			}
			ctx.Locals(CtxKeyUserID, out.UserID)
			ctx.Locals(CtxKeySessionID, out.SessionID)

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

func findSessionf(uc UsecasesForRequireAuthorizedSession, token string) (findSession.Out, error) {
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

type OutJwt struct {
	UserID    string
	SessionID string
}

type JwtParser interface {
	Parse(token string) (OutJwt, error)
}

func parseJwt(ucjwt JwtParser, token string) (OutJwt, error) {
	out, err := ucjwt.Parse(token)
	if err != nil {
		return OutJwt{}, fiber.ErrUnauthorized
	}

	return out, nil
}
