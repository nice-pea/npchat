package router

import (
	"errors"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
)

type ErrCode interface {
	ErrCode() string
}

// errCode по ошибке определяет код
func errCode(err error) string {
	var errWithCode ErrCode
	if errors.As(err, &errWithCode) {
		return errWithCode.ErrCode()
	}

	switch {
	case errors.Is(err, middleware.ErrUnauthorized):
		return ErrCodeInvalidAuthorizationHeader
	case errors.Is(err, ErrJsonMarshalResponseData):
		return ErrCodeUnmarshalJSONResponseData
	case errors.Is(err, middleware.ErrUnknownRequestID):
		return ErrCodeInvalidXRequestIDHeader
	case errors.Is(err, middleware.ErrUnsupportedAcceptValue):
		return ErrCodeUnsupportedAccept
	case errors.Is(err, middleware.ErrUnsupportedContentTypeValue):
		return ErrCodeUnsupportedContentType
	}

	return ErrCodeUnknown
}

const (
	ErrCodeUnknown                    = ""
	ErrCodeInvalidAuthorizationHeader = "INVALID_AUTHORIZATION_HEADER"
	ErrCodeInvalidXRequestIDHeader    = "INVALID_X_REQUEST_ID_HEADER"
	ErrCodeUnsupportedAccept          = "UNSUPPORTED_CONTENT_TYPE"
	ErrCodeUnsupportedContentType     = "UNSUPPORTED_ACCEPTED_CONTENT_TYPE"
	ErrCodeUnmarshalJSONResponseData  = "UNMARSHAL_JSON_RESPONSE_DATA"
)
