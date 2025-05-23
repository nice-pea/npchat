package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// DeleteMember регистрирует обработчик, позволяющий удалить участника из чата.
// Доступен только авторизованным пользователям, которые являются главными администраторами чата.
// Метод: DELETE /chats/{chatID}/members
func DeleteMember(router http2.Router) {
	// requestBody описывает структуру тела запроса для удаления участника из чата.
	type requestBody struct {
		UserID string `json:"user_id"` // ID пользователя, которого требуется удалить из чата
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

			// Формируем входные данные для сервиса удаления участника.
			// SubjectUserID - ID пользователя, выполняющего удаление (должен быть главным администратором).
			// ChatID - ID чата, из которого удаляется участник (берётся из параметра пути).
			// UserID - ID пользователя, которого требуется удалить (берётся из тела запроса)
			input := service.DeleteMemberInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        http2.PathStr(context, "chatID"),
				UserID:        rb.UserID,
			}

			// Вызываем сервис удаления участника из чата.
			return nil, context.Services().Members().DeleteMember(input)
		})
}
