package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// MyInvitations регистрирует обработчик, позволяющий получить список приглашений пользователя.
// Доступен только авторизованным пользователям.
//
// Метод: GET /invitations
func MyInvitations(router http2.Router) {
	router.HandleFunc(
		"GET /invitations",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			input := service.ReceivedInvitationsInput{
				SubjectUserID: context.Session().UserID,
				UserID:        context.Session().UserID,
			}

			return context.Services().Invitations().ReceivedInvitations(input)
		})
}
