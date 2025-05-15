package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// ChatInvitations регистрирует обработчик, позволяющий Получить список приглашений в чат
func ChatInvitations(router http2.Router) {
	router.HandleFunc(
		"GET /chats/{chatID}/invitations",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			input := service.ChatInvitationsInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        http2.PathStr(context, "chatID"),
			}
			return context.Services().Invitations().ChatInvitations(input)
		})
}
