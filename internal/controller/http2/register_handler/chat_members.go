package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	chatMembers "github.com/nice-pea/npchat/internal/usecases/chats/chat_members"
)

// ChatMembers регистрирует обработчик, позволяющий получить список участников чата.
// Доступен только авторизованным пользователям.
//
// Метод: GET /chats/{chatID}/members
func ChatMembers(router *fiber.App, uc UsecasesForChatMembers, jwtParser middleware.JwtParser) {
	router.Get(
		"/chats/:chatID/members",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
		func(ctx *fiber.Ctx) error {
			input := chatMembers.In{
				SubjectID: UserID(ctx),
				ChatID:    ParamsUUID(ctx, "chatID"),
			}

			out, err := uc.ChatMembers(input)
			if err != nil {
				return err
			}

			return ctx.JSON(out)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
	)
}

// UsecasesForChatMembers определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForChatMembers interface {
	ChatMembers(chatMembers.In) (chatMembers.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
