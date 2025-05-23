package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// LeaveChat регистрирует обработчик, позволяющий пользователю покинуть чат.
// Доступен только авторизованным пользователям.
// Метод: POST /chats/{chatID}/leave
func LeaveChat(router http2.Router) {
	router.HandleFunc(
		"POST /chats/{chatID}/leave",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			// Формируем входные данные для выхода пользователя из чата.
			// SubjectUserID - ID пользователя, который выходит из чата.
			// ChatID - ID чата, который пользователь покидает (берётся из параметра пути).
			input := service.LeaveChatInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        http2.PathStr(context, "chatID"),
			}

			// Вызываем сервис выхода из чата.
			return nil, context.Services().Members().LeaveChat(input)
		})
}
