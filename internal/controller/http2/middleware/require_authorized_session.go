package middleware

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/service"
)

// ErrUnauthorized запрос не содержит действительный токен авторизации
var ErrUnauthorized = errors.New("unauthorized. Please, use token in header: Authorization: Bearer <token>")

// RequireAuthorizedSession требует авторизованную сессии
func RequireAuthorizedSession(next http2.HandlerFuncRW) http2.HandlerFuncRW {
	return func(context http2.RWContext) (any, error) {
		var err error
		session, err := getSession(context)
		if err != nil {
			return nil, fmt.Errorf("getSession: %w", err)
		}
		if session.ID == uuid.Nil {
			return nil, ErrUnauthorized
		}
		context.SetSession(session)

		return next(context)
	}
}

func getSession(ctx http2.Context) (sessionn.Session, error) {
	header := ctx.Request().Header.Get("Authorization")
	token, _ := strings.CutPrefix(header, "Bearer ")
	if token == "" {
		return sessionn.Session{}, nil
	}

	sessions, err := ctx.Services().Sessions().Find(service.SessionsFindIn{
		Token: token,
	})
	if err != nil {
		return sessionn.Session{}, fmt.Errorf("c.sessions.Find: %w", err)
	}
	if len(sessions) != 1 {
		return sessionn.Session{}, nil
	}

	return sessions[0], nil
}
