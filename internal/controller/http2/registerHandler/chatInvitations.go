package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// ChatInvitations регистрирует обработчик, позволяющий получить список приглашений в определённый чат.
// Доступен только авторизованным пользователям.
//
// Метод: GET /chats/{chatID}/invitations
func ChatInvitations(router http2.Router) {
	router.HandleFunc(
		"GET /chats/{chatID}/invitations",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			input := service.ChatInvitationsIn{
				SubjectID: context.Session().UserID,
				ChatID:    http2.PathStr(context, "chatID"),
			}

			return context.Services().Invitations().ChatInvitations(input)
		})
}
