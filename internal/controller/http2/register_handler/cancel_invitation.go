package register_handler

import (
	"github.com/nice-pea/npchat/internal/controller/http2"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// CancelInvitation регистрирует обработчик, позволяющий отменить приглашение в чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /invitations/{invitationID}/cancel
func CancelInvitation(router http2.Router) {
	router.HandleFunc(
		"POST /invitations/{invitationID}/cancel",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			input := service.CancelInvitationIn{
				SubjectID:    context.Session().UserID,
				InvitationID: http2.PathUUID(context, "invitationID"),
			}

			return nil, context.Services().Chats().CancelInvitation(input)
		})
}
