package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// MyInvitations регистрирует обработчик, позволяющий получить список приглашений пользователя.
// Доступен только авторизованным пользователям.
//
// Метод: GET /invitations
func MyInvitations(router *fiber.App, ss Services) {
	router.Get(
		"/invitations",
		func(context *fiber.Ctx) error {
			input := service.ReceivedInvitationsIn{
				SubjectID: Session(context).UserID,
			}

			out, err := ss.Chats().ReceivedInvitations(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequareAuthoruzation(ss.Sessions()),
	)
}
