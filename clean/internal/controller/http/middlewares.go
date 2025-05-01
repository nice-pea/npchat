package http

import (
	"errors"
)

var ErrUnauthorized = errors.New("unauthorized. Please, use token in header: Authorization: Bearer <token>")

func requireAuthorizedSession(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		if context.session.ID == "" {
			return nil, ErrUnauthorized
		}

		return next(context)
	}
}

var ErrUnknownRequestID = errors.New("unknown request ID. Please, use X-Request-ID header")

func requireRequestID(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		if context.requestID == "" {
			return nil, ErrUnknownRequestID
		}

		return next(context)
	}
}

var ErrUnsupportedAcceptedContentType = errors.New("unsupported Accept header value. Please, use Accept: application/json header")

func requireAcceptJson(next HandlerFunc) HandlerFunc {
	return func(context Context) (any, error) {
		if context.request.Header.Get("Accept") != "application/json" {
			return nil, ErrUnsupportedAcceptedContentType
		}

		return next(context)
	}
}
