package middleware

import (
	"errors"
	"strings"

	"github.com/nice-pea/npchat/internal/controller/http2"
)

// ErrUnsupportedAcceptValue клиент не принимает JSON в качестве ответа
var ErrUnsupportedAcceptValue = errors.New("unsupported Accept header value. Please, use Accept: application/json header")

// RequireAcceptJson требует поддержку json как типа контента, который ожидание клиент
func RequireAcceptJson(next http2.HandlerFuncRW) http2.HandlerFuncRW {
	return func(context http2.RWContext) (any, error) {
		accept := context.Request().Header.Get("Accept")
		if strings.Contains(accept, "application/json") || accept == "" || accept == "*/*" {
			return next(context)
		}

		return nil, ErrUnsupportedAcceptValue
	}
}
