package register_handler

import (
	"github.com/nice-pea/npchat/internal/controller/http2"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// AcceptInvitation регистрирует обработчик, позволяющий принять приглашение в чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /invitations/{invitationID}/accept
func AcceptInvitation(router http2.Router) {
	router.HandleFunc(
		"POST /invitations/{invitationID}/accept",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			input := service.AcceptInvitationIn{
				SubjectID:    context.Session().UserID,
				InvitationID: http2.PathUUID(context, "invitationID"),
			}

			return nil, context.Services().Chats().AcceptInvitation(input)
		})
}
