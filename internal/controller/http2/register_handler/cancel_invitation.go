package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	cancelInvitation "github.com/nice-pea/npchat/internal/usecases/chats/cancel_invitation"
)

// CancelInvitation регистрирует обработчик, позволяющий отменить приглашение в чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /invitations/{invitationID}/cancel
func CancelInvitation(router *fiber.App, uc UsecasesForCancelInvitation, jparser middleware.JwtParser) {
	router.Post(
		"/invitations/:invitationID/cancel",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jparser),
		func(context *fiber.Ctx) error {
			input := cancelInvitation.In{
				SubjectID:    UserID(context),
				InvitationID: ParamsUUID(context, "invitationID"),
			}

			out, err := uc.CancelInvitation(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// UsecasesForCancelInvitation определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForCancelInvitation interface {
	CancelInvitation(cancelInvitation.In) (cancelInvitation.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
