package register_handler

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
func MyInvitations(router *fiber.App, uc UsecasesForMyInvitations) {
	router.Get(
		"/invitations",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc),
		func(context *fiber.Ctx) error {
			input := receivedInvitations.In{
				SubjectID: Session(context).UserID,
			}

			out, err := uc.ReceivedInvitations(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// UsecasesForMyInvitations определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForMyInvitations interface {
	ReceivedInvitations(receivedInvitations.In) (receivedInvitations.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
