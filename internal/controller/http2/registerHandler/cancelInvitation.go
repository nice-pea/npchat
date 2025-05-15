package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Отменить приглашение в чат
func CancelInvitation(router http2.Router) {
	router.HandleFunc(
		"POST /invitations/{invitationID}/cancel",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			input := service.CancelInvitationInput{
				SubjectUserID: context.Session().UserID,
				UserID:        "",
				ChatID:        "",
			}
			return nil, context.Services().Invitations().CancelInvitation()
		})
}
