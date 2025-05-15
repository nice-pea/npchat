package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

// Получить список приглашений в чат
func RegisterChatInvitationsHandler(router http2.Router) {
	router.HandleFunc(
		"GET /chats/{chatID}/invitations",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
