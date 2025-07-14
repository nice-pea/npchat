package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// MyChats регистрирует HTTP-обработчик для получения списка чатов пользователя.
// Данный обработчик доступен только авторизованным пользователям.
//
// Метод: GET /chats
func MyChats(router *fiber.App, ss services) {
	router.Get(
		"/chats",
		func(context *fiber.Ctx) error {
			input := service.WhichParticipateIn{
				SubjectID: Session(context).UserID,
				UserID:    Session(context).UserID,
			}

			out, err := ss.Chats().WhichParticipate(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequareAuthoruzation(ss.Sessions()),
	)
}
