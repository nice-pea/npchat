package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Получить список моих приглашений
func MyInvitations(router http2.Router) {
	router.HandleFunc(
		"GET /invitations",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			input := service.UserInvitationsInput{
				SubjectUserID: context.Session().UserID,
				UserID:        context.Session().UserID,
			}
			return context.Services().Invitations().UserInvitations(input)
		})
}
