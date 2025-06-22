package middleware

import (
	"errors"

	"github.com/nice-pea/npchat/internal/controller/http2"
)

// ErrUnknownRequestID запрос не содержит идентификатор запроса
var ErrUnknownRequestID = errors.New("unknown request ID. Please, use X-Request-ID header")

// RequireRequestID требует наличие идентификатора запроса
func RequireRequestID(next http2.HandlerFuncRW) http2.HandlerFuncRW {
	return func(context http2.RWContext) (any, error) {
		context.SetRequestID(context.Request().Header.Get("X-Request-ID"))
		if context.RequestID() == "" {
			return nil, ErrUnknownRequestID
		}

		return next(context)
	}
}
