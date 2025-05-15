package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

// Обновить название чата
func RegisterUpdateChatNameHandler(router http2.Router) {
	router.HandleFunc(
		"PUT /chats/{chatID}/name",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
