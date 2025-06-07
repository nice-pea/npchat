package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// CreateChat регистрирует обработчик, позволяющий создать новый чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /chats
func CreateChat(router http2.Router) {
	// Тело запроса для создания чата.
	type requestBody struct {
		Name string `json:"name"`
	}
	router.HandleFunc(
		"POST /chats",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.CreateChatIn{
				ChiefUserID: context.Session().UserID,
				Name:        rb.Name,
			}

			return context.Services().Chats().CreateChat(input)
		})
}
