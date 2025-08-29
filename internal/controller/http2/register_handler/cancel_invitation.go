package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// CancelInvitation регистрирует обработчик, позволяющий отменить приглашение в чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /invitations/{invitationID}/cancel
func CancelInvitation(router *fiber.App, ss Services) {
	router.Post(
		"/invitations/:invitationID/cancel",
		func(context *fiber.Ctx) error {
			input := service.CancelInvitationIn{
				SubjectID:    Session(context).UserID,
				InvitationID: ParamsUUID(context, "invitationID"),
			}

			return ss.Chats().CancelInvitation(input)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(ss),
	)
}
