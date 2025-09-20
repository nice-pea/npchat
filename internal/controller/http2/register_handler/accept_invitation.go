package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	acceptInvitation "github.com/nice-pea/npchat/internal/usecases/chats/accept_invitation"
)

// AcceptInvitation регистрирует обработчик, позволяющий принять приглашение в чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /invitations/{invitationID}/accept
func AcceptInvitation(router *fiber.App, uc UsecasesForAcceptInvitation, jwtParser middleware.JwtParser) {
	router.Post(
		"/invitations/:invitationID/accept",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
		func(ctx *fiber.Ctx) error {
			input := acceptInvitation.In{
				SubjectID:    UserID(ctx),
				InvitationID: ParamsUUID(ctx, "invitationID"),
			}

			out, err := uc.AcceptInvitation(input)
			if err != nil {
				return err
			}

			return ctx.JSON(out)
		},
	)
}

// UsecasesForAcceptInvitation определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForAcceptInvitation interface {
	AcceptInvitation(acceptInvitation.In) (acceptInvitation.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
