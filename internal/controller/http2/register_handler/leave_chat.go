package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// LeaveChat регистрирует обработчик, позволяющий пользователю покинуть чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /chats/{chatID}/leave
func LeaveChat(router *fiber.App, ss Services) {
	router.Post(
		"/chats/:chatID/leave",
		func(context *fiber.Ctx) error {
			input := service.LeaveChatIn{
				SubjectID: Session(context).UserID,
				ChatID:    ParamsUUID(context, "chatID"),
			}

			return ss.Chats().LeaveChat(input)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(ss.Sessions()),
	)
}
