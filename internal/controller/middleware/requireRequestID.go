package middleware

import (
	"errors"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
)

// ErrUnknownRequestID запрос не содержит идентификатор запроса
var ErrUnknownRequestID = errors.New("unknown request ID. Please, use X-Request-ID header")

// RequireRequestID требует наличие идентификатора запроса
func RequireRequestID(next http2.MiddlewareFunc) http2.MiddlewareFunc {
	return func(context http2.MutContext) (any, error) {
		context.SetRequestID(context.Request().Header.Get("X-Request-ID"))
		if context.RequestID() == "" {
			return nil, ErrUnknownRequestID
		}

		return next(context)
	}
}
