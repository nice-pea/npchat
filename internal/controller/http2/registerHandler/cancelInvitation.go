package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// CancelInvitation регистрирует обработчик, позволяющий отменить приглашение в чат.
// Доступен только авторизованным пользователям.
// Метод: POST /invitations/{invitationID}/cancel
func CancelInvitation(router http2.Router) {
	router.HandleFunc(
		"POST /invitations/{invitationID}/cancel",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			// Формируем входные данные для отмены приглашения.
			// SubjectUserID - ID пользователя, выполняющего отмену.
			// InvitationID - ID приглашения, которое требуется отменить.
			input := service.CancelInvitationInput{
				SubjectUserID: context.Session().UserID,
				InvitationID:  http2.PathStr(context, "invitationID"),
			}

			// Вызываем сервис отмены приглашения.
			return nil, context.Services().Invitations().CancelInvitation(input)
		})
}
