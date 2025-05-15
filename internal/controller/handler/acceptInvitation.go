package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

// Принять приглашение в чат
func RegisterAcceptInvitationHandler(router http2.Router) {
	router.HandleFunc(
		"POST /invitations/{invitationID}/accept",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
