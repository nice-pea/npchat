package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	receivedInvitations "github.com/nice-pea/npchat/internal/usecases/chats/received_invitations"
)

// MyInvitations регистрирует обработчик, позволяющий получить список приглашений пользователя.
// Доступен только авторизованным пользователям.
//
// Метод: GET /invitations
func MyInvitations(router *fiber.App, uc UsecasesForMyInvitations, jwtParser middleware.JwtParser) {
	router.Get(
		"/invitations",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
		func(ctx *fiber.Ctx) error {
			input := receivedInvitations.In{
				SubjectID: UserID(ctx),
			}

			out, err := uc.ReceivedInvitations(input)
			if err != nil {
				return err
			}

			return ctx.JSON(out)
		},
	)
}

// UsecasesForMyInvitations определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForMyInvitations interface {
	ReceivedInvitations(receivedInvitations.In) (receivedInvitations.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
