package register_handler

import (
	"github.com/nice-pea/npchat/internal/controller/http2"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// ChatMembers регистрирует обработчик, позволяющий получить список участников чата.
// Доступен только авторизованным пользователям.
//
// Метод: GET /chats/{chatID}/members
func ChatMembers(router http2.Router) {
	router.HandleFunc(
		"GET /chats/{chatID}/members",
		middleware.ClientAuthChain, // Цепочка middleware для клиентских запросов с аутентификацией
		func(context http2.Context) (any, error) {
			input := service.ChatMembersIn{
				SubjectID: context.Session().UserID,
				ChatID:    http2.PathUUID(context, "chatID"),
			}

			return context.Services().Chats().ChatMembers(input)
		})
}
