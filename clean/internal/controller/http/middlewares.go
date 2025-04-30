package http

import (
	"errors"
)

var ErrUnauthorized = errors.New("unauthorized. Please, use token in header: Authorization: Bearer <token>")

func requireAuthorizedSession(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		println("requireAuthorizedSession")
		if context.session.ID == "" {
			return nil, ErrUnauthorized
		}

		return next(context)
	}
}

var ErrUnknownRequestID = errors.New("unknown request ID. Please, use X-Request-ID header")

func requireRequestID(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		println("requireRequestID")
		if context.requestID == "" {
			return nil, ErrUnknownRequestID
		}

		return next(context)
	}
}

var ErrOnlyJsonSupported = errors.New("only JSON format is supported. Please, use Accept: application/json header")

func requireAcceptJson(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		println("requireAcceptJson")
		if context.request.Header.Get("Accept") != "application/json" {
			return nil, ErrOnlyJsonSupported
		}

		return next(context)
	}
}
