package registerHandler

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
	// requestBody описывает структуру тела запроса для отправки приглашения.
	type requestBody struct {
		ChatID string `json:"chat_id"` // ID чата, в который отправляется приглашение
		UserID string `json:"user_id"` // ID пользователя, которого приглашают в чат
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

			// Формируем входные данные для сервиса отправки приглашения.
			// SubjectUserID - ID пользователя, отправляющего приглашение (берётся из сессии)
			// ChatID - ID чата, в который отправляется приглашение.
			// UserID - ID пользователя, которого приглашают.
			input := service.SendInvitationInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        rb.ChatID,
				UserID:        rb.UserID,
			}

			// Вызываем сервис отправки приглашения и возвращаем результат.
			return context.Services().Invitations().SendInvitation(input)
		})
}
