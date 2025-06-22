package register_handler

import (
	"github.com/nice-pea/npchat/internal/controller/http2"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// LeaveChat регистрирует обработчик, позволяющий пользователю покинуть чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /chats/{chatID}/leave
func LeaveChat(router http2.Router) {
	router.HandleFunc(
		"POST /chats/{chatID}/leave",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			input := service.LeaveChatIn{
				SubjectID: context.Session().UserID,
				ChatID:    http2.PathUUID(context, "chatID"),
			}

			return nil, context.Services().Chats().LeaveChat(input)
		})
}
