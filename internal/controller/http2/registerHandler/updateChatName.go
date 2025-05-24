package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// UpdateChatName регистрирует обработчик, позволяющий обновить название чата.
// Доступен только авторизованным пользователям, которые являются главными администраторами чата.
//
// Метод: PUT /chats/{chatID}/name
func UpdateChatName(router http2.Router) {
	// requestBody описывает структуру тела запроса для обновления названия чата.
	type requestBody struct {
		NewName string `json:"new_name"` // Новое название чата
	}
	router.HandleFunc(
		"PUT /chats/{chatID}/name",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			// Формируем входные данные для сервиса обновления названия чата.
			// SubjectUserID - ID пользователя, выполняющего обновление.
			// ChatID - ID чата, название которого обновляется (берётся из параметра пути).
			// NewName - новое название чата, полученное из запроса.
			input := service.UpdateNameInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        http2.PathStr(context, "chatID"),
				NewName:       rb.NewName,
			}

			// Вызываем сервис обновления названия чата и возвращаем результат.
			return context.Services().Chats().UpdateName(input)
		})
}
