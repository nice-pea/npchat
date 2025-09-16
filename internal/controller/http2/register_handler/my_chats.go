package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	myChats "github.com/nice-pea/npchat/internal/usecases/chats/my_chats"
)

// MyChats регистрирует HTTP-обработчик для получения списка чатов пользователя.
// Данный обработчик доступен только авторизованным пользователям.
//
// Метод: GET /chats
func MyChats(router *fiber.App, uc UsecasesForMyChats, jparser middleware.JwtParser) {
	router.Get(
		"/chats",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jparser),
		func(context *fiber.Ctx) error {
			input := myChats.In{
				SubjectID: UserID(context),
				UserID:    UserID(context),
			}

			out, err := uc.MyChats(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// UsecasesForMyChats определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForMyChats interface {
	MyChats(myChats.In) (myChats.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
