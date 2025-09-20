package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	leaveChat "github.com/nice-pea/npchat/internal/usecases/chats/leave_chat"
)

// LeaveChat регистрирует обработчик, позволяющий пользователю покинуть чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /chats/{chatID}/leave
func LeaveChat(router *fiber.App, uc UsecasesForLeaveChat, jwtParser middleware.JwtParser) {
	router.Post(
		"/chats/:chatID/leave",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
		func(context *fiber.Ctx) error {
			input := leaveChat.In{
				SubjectID: UserID(context),
				ChatID:    ParamsUUID(context, "chatID"),
			}

			out, err := uc.LeaveChat(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// UsecasesForLeaveChat определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForLeaveChat interface {
	LeaveChat(leaveChat.In) (leaveChat.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
