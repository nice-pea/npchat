package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// ChatInvitations регистрирует обработчик, позволяющий получить список приглашений в определённый чат.
// Доступен только авторизованным пользователям.
//
// Метод: GET /chats/{chatID}/invitations
func ChatInvitations(router *fiber.App, ss Services) {
	router.Get(
		"/chats/:chatID/invitations",
		func(context *fiber.Ctx) error {
			input := service.ChatInvitationsIn{
				SubjectID: Session(context).UserID,
				ChatID:    ParamsUUID(context, "chatID"),
			}

			out, err := ss.Chats().ChatInvitations(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(ss),
	)
}
