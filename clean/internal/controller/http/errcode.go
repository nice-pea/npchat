package http

import (
	"errors"
)

func errCode(err error) string {
	var errWithCode interface {
		ErrCode() string
		Error() string
	}
	if errors.As(err, &errWithCode) {
		return errWithCode.ErrCode()
	}

	switch {
	case errors.Is(err, ErrUnauthorized):
		return ErrCodeInvalidAuthorizationHeader
	case errors.Is(err, ErrJsonMarshalResponseData):
		return ErrCodeUnmarshalJSONResponseData
	case errors.Is(err, ErrUnknownRequestID):
		return ErrCodeInvalidXRequestIDHeader
	case errors.Is(err, ErrUnsupportedAcceptedContentType):
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
