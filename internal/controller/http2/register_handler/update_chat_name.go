package register_handler

import (
	"github.com/nice-pea/npchat/internal/controller/http2"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// UpdateChatName регистрирует обработчик, позволяющий обновить название чата.
// Доступен только авторизованным пользователям, которые являются главными администраторами чата.
//
// Метод: PUT /chats/{chatID}/name
func UpdateChatName(router http2.Router) {
	// Тело запроса для обновления названия чата.
	type requestBody struct {
		NewName string `json:"new_name"`
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

			input := service.UpdateNameIn{
				SubjectID: context.Session().UserID,
				ChatID:    http2.PathUUID(context, "chatID"),
				NewName:   rb.NewName,
			}

			return context.Services().Chats().UpdateName(input)
		})
}
