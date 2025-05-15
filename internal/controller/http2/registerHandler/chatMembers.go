package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
)

// Получить список участников чата
func ChatMembers(router http2.Router) {
	router.HandleFunc(
		"GET /chats/{chatID}/members",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
