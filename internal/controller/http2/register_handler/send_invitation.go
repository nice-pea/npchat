package register_handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// SendInvitation регистрирует обработчик, позволяющий отправить приглашение в чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /invitations
func SendInvitation(router http2.Router) {
	// Тело запроса для отправки приглашения.
	type requestBody struct {
		ChatID string `json:"chat_id"`
		UserID string `json:"user_id"`
	}
	router.HandleFunc(
		"POST /invitations",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.SendInvitationIn{
				SubjectID: context.Session().UserID,
				ChatID:    rb.ChatID,
				UserID:    rb.UserID,
			}

			return context.Services().Chats().SendInvitation(input)
		})
}
