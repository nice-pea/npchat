package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// ChatMembers регистрирует обработчик, позволяющий получить список участников чата.
// Доступен только авторизованным пользователям.
// Метод: GET /chats/{chatID}/members
func ChatMembers(router http2.Router) {
	router.HandleFunc(
		"GET /chats/{chatID}/members",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			// Формируем входные данные для получения участников чата.
			// SubjectUserID - ID пользователя, выполняющего запрос.
			// ChatID - ID чата, для которого запрашиваются участники.
			input := service.ChatMembersInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        http2.PathStr(context, "chatID"),
			}

			// Вызываем сервис получения участников чата и возвращаем результат.
			return context.Services().Members().ChatMembers(input)
		})
}
