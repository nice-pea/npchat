package register_handler

import (
	"github.com/nice-pea/npchat/internal/controller/http2"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
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
			input := service.ReceivedInvitationsIn{
				SubjectID: context.Session().UserID,
			}

			return context.Services().Chats().ReceivedInvitations(input)
		})
}
