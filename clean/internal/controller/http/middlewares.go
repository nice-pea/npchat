package http

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// ErrUnauthorized запрос не содержит действительный токен авторизации
var ErrUnauthorized = errors.New("unauthorized. Please, use token in header: Authorization: Bearer <token>")

// requireAuthorizedSession требует авторизованную сессии
func (c *Controller) requireAuthorizedSession(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		var err error
		if context.session, err = c.getSession(context.request); err != nil {
			return nil, fmt.Errorf("getSession: %w", err)
		}
		if context.session.ID == "" {
			return nil, ErrUnauthorized
		}

		return next(context)
	}
}

func (c *Controller) getSession(r *http.Request) (domain.Session, error) {
	header := r.Header.Get("Authorization")
	token, _ := strings.CutPrefix(header, "Bearer ")
	if token == "" {
		return domain.Session{}, nil
	}

	sessions, err := c.sessions.Find(service.SessionsFindInput{
		Token: token,
	})
	if err != nil {
		return domain.Session{}, fmt.Errorf("c.sessions.Find: %w", err)
	}
	if len(sessions) != 1 {
		return domain.Session{}, nil
	}

	return sessions[0], nil
}

// ErrUnknownRequestID запрос не содержит идентификатор запроса
var ErrUnknownRequestID = errors.New("unknown request ID. Please, use X-Request-ID header")

// requireRequestID требует наличие идентификатора запроса
func (c *Controller) requireRequestID(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		context.requestID = context.request.Header.Get("X-Request-ID")
		if context.requestID == "" {
			return nil, ErrUnknownRequestID
		}

		return next(context)
	}
}

// ErrUnsupportedAcceptedContentType клиент не принимает JSON в качестве ответа
var ErrUnsupportedAcceptedContentType = errors.New("unsupported Accept header value. Please, use Accept: application/json header")

// requireAcceptJson требует поддержку json как типа контента, который ожидание клиент
func requireAcceptJson(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		if context.request.Header.Get("Accept") != "application/json" {
			return nil, ErrUnsupportedAcceptedContentType
		}

		return next(context)
	}
}
