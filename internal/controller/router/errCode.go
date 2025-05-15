package router

import (
	"errors"

	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

type ErrCode interface {
	ErrCode() string
}

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
	case errors.Is(err, middleware.ErrUnsupportedAcceptedContentType):
		return ErrCodeUnsupportedAcceptedContentType
	}

	return ErrCodeUnknown
}

const (
	ErrCodeUnknown                        = ""
	ErrCodeInvalidAuthorizationHeader     = "INVALID_AUTHORIZATION_HEADER"
	ErrCodeInvalidXRequestIDHeader        = "INVALID_X_REQUEST_ID_HEADER"
	ErrCodeUnsupportedAcceptedContentType = "UNSUPPORTED_ACCEPTED_CONTENT_TYPE"
	ErrCodeUnmarshalJSONResponseData      = "UNMARSHAL_JSON_RESPONSE_DATA"
)
