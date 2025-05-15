package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

// Отменить приглашение в чат
func RegisterCancelInvitationHandler(router http2.Router) {
	router.HandleFunc(
		"POST /invitations/{invitationID}/cancel",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
