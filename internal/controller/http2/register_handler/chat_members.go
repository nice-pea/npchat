package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// ChatMembers регистрирует обработчик, позволяющий получить список участников чата.
// Доступен только авторизованным пользователям.
//
// Метод: GET /chats/{chatID}/members
func ChatMembers(router *fiber.App, ss services) {
	router.Get(
		"/chats/:chatID/members",
		func(context *fiber.Ctx) error {
			input := service.ChatMembersIn{
				SubjectID: Session(context).UserID,
				ChatID:    ParamsUUID(context, "chatID"),
			}

			out, err := ss.Chats().ChatMembers(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequareAuthoruzation(ss.Sessions()),
	)
}
