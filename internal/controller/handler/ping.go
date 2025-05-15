package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

func RegisterPingHandler(router http2.Router) {
	router.HandleFunc(
		"/ping",
		middleware.EmptyChain,
		func(context http2.Context) (any, error) {
			return "pong", nil
		})
}
