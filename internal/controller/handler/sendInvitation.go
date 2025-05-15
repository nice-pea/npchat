package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

// Отправить приглашение в чат
func RegisterSendInvitationHandler(router http2.Router) {
	router.HandleFunc(
		"POST /invitations",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
