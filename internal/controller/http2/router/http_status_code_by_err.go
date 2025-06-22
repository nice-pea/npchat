package router

import (
	"errors"
	"net/http"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
)

func httpStatusCodeByErr(err error) int {
	if errors.Is(err, middleware.ErrUnauthorized) {
		return http.StatusUnauthorized
	}
	if errors.Is(err, ErrJsonMarshalResponseData) {
		return http.StatusInternalServerError
	}
	if errors.Is(err, ErrWriteResponseBytes) {
		return http.StatusInternalServerError
	}

	return http.StatusBadRequest
}
