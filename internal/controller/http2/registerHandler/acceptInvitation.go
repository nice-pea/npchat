package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Принять приглашение в чат
func AcceptInvitation(router http2.Router) {
	router.HandleFunc(
		"POST /invitations/{invitationID}/accept",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			input := service.AcceptInvitationInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        "",
			}
			return nil, context.Services().Invitations().AcceptInvitation()
		})
}
