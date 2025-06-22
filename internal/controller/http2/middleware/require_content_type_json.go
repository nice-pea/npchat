package middleware

import (
	"errors"

	"github.com/nice-pea/npchat/internal/controller/http2"
)

// ErrUnsupportedContentTypeValue клиент не принимает JSON в качестве ответа
var ErrUnsupportedContentTypeValue = errors.New("unsupported ContentType header value. Please, use ContentType: application/json header")

// RequireContentTypeJson требует поддержку json как типа контента, который отправляет клиент
func RequireContentTypeJson(next http2.HandlerFuncRW) http2.HandlerFuncRW {
	return func(context http2.RWContext) (any, error) {
		if context.Request().Header.Get("Content-Type") != "application/json" {
			return nil, ErrUnsupportedContentTypeValue
		}

		return next(context)
	}
}
