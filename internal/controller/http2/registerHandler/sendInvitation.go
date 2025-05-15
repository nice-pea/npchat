package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// SendInvitation регистрирует обработчик, позволяющий Отправить приглашение в чат
func SendInvitation(router http2.Router) {
	type requestBody struct {
		ChatID string `json:"chat_id"`
		UserID string `json:"user_id"`
	}
	router.HandleFunc(
		"POST /invitations",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.SendInvitationInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        rb.ChatID,
				UserID:        rb.UserID,
			}

			return context.Services().Invitations().SendInvitation(input)
		})
}
