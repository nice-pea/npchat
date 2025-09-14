package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	chatInvitations "github.com/nice-pea/npchat/internal/usecases/chats/chat_invitations"
)

// ChatInvitations регистрирует обработчик, позволяющий получить список приглашений в определённый чат.
// Доступен только авторизованным пользователям.
//
// Метод: GET /chats/{chatID}/invitations
func ChatInvitations(router *fiber.App, uc UsecasesForChatInvitations, jparser middleware.JWTParser) {
	router.Get(
		"/chats/:chatID/invitations",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jparser),
		func(context *fiber.Ctx) error {
			input := chatInvitations.In{
				SubjectID: Session(context).UserID,
				ChatID:    ParamsUUID(context, "chatID"),
			}

			out, err := uc.ChatInvitations(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jparser),
	)
}

// UsecasesForChatInvitations определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForChatInvitations interface {
	ChatInvitations(chatInvitations.In) (chatInvitations.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
