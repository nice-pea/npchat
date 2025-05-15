package middleware

import (
	"errors"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
)

// ErrUnsupportedAcceptValue клиент не принимает JSON в качестве ответа
var ErrUnsupportedAcceptValue = errors.New("unsupported Accept header value. Please, use Accept: application/json header")

// RequireAcceptJson требует поддержку json как типа контента, который ожидание клиент
func RequireAcceptJson(next http2.HandlerFuncRW) http2.HandlerFuncRW {
	return func(context http2.RWContext) (any, error) {
		if context.Request().Header.Get("Accept") != "application/json" {
			return nil, ErrUnsupportedAcceptValue
		}

		return next(context)
	}
}
