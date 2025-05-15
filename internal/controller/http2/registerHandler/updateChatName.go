package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
)

// Обновить название чата
func UpdateChatName(router http2.Router) {
	router.HandleFunc(
		"PUT /chats/{chatID}/name",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
