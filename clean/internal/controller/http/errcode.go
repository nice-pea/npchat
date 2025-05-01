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

	if errors.Is(err, ErrUnauthorized) {
		return ErrCodeInvalidAuthorizationHeader
	}

	return ErrCodeUnknown
}

const (
	ErrCodeUnknown                    = ""
	ErrCodeInvalidAuthorizationHeader = "INVALID_AUTHORIZATION_HEADER"
)
