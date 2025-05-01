package http

import (
	"errors"
)

// ErrUnauthorized запрос не содержит действительный токен авторизации
var ErrUnauthorized = errors.New("unauthorized. Please, use token in header: Authorization: Bearer <token>")

// requireAuthorizedSession требует авторизованную сессии
func requireAuthorizedSession(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		if context.session.ID == "" {
			return nil, ErrUnauthorized
		}

		return next(context)
	}
}

// ErrUnknownRequestID запрос не содержит идентификатор запроса
var ErrUnknownRequestID = errors.New("unknown request ID. Please, use X-Request-ID header")

// requireRequestID требует наличие идентификатора запроса
func requireRequestID(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
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
