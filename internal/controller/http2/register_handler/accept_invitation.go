package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// AcceptInvitation регистрирует обработчик, позволяющий принять приглашение в чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /invitations/{invitationID}/accept
func AcceptInvitation(router *fiber.App, ss Services) {
	router.Post(
		"/invitations/:invitationID/accept",
		func(context *fiber.Ctx) error {
			input := service.AcceptInvitationIn{
				SubjectID:    Session(context).UserID,
				InvitationID: ParamsUUID(context, "invitationID"),
			}

			return ss.Chats().AcceptInvitation(input)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(ss),
	)
}
