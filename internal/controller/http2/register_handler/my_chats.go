package register_handler

import (
	"github.com/nice-pea/npchat/internal/controller/http2"
	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// MyChats регистрирует HTTP-обработчик для получения списка чатов пользователя.
// Данный обработчик доступен только авторизованным пользователям.
//
// Метод: GET /chats
func MyChats(router http2.Router) {
	router.HandleFunc(
		"GET /chats",
		middleware.ClientAuthChain, // Цепочка middleware, проверяющая авторизацию пользователя
		func(context http2.Context) (any, error) {
			input := service.WhichParticipateIn{
				SubjectID: context.Session().UserID,
				UserID:    context.Session().UserID,
			}

			return context.Services().Chats().WhichParticipate(input)
		})
}
