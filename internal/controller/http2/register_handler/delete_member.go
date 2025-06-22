package register_handler

import (
	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// DeleteMember регистрирует обработчик, позволяющий удалить участника из чата.
// Доступен только авторизованным пользователям, которые являются главными администраторами чата.
//
// Метод: DELETE /chats/{chatID}/members
func DeleteMember(router http2.Router) {
	// Тело запроса для удаления участника из чата.
	type requestBody struct {
		UserID uuid.UUID `json:"user_id"`
	}
	router.HandleFunc(
		"DELETE /chats/{chatID}/members",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.DeleteMemberIn{
				SubjectID: context.Session().UserID,
				ChatID:    http2.PathUUID(context, "chatID"),
				UserID:    rb.UserID,
			}

			return nil, context.Services().Chats().DeleteMember(input)
		})
}
