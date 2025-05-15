package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

// Получить список моих приглашений
func RegisterMyInvitationsHandler(router http2.Router) {
	router.HandleFunc(
		"GET /invitations",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
