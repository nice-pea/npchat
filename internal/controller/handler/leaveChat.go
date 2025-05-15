package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

// Покинуть чат
func RegisterLeaveChatHandler(router http2.Router) {
	router.HandleFunc(
		"POST /chats/{chatID}/leave",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
