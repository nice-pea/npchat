package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	createChat "github.com/nice-pea/npchat/internal/service/chats/create_chat"
)

// CreateChat регистрирует обработчик, позволяющий создать новый чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /chats
func CreateChat(router *fiber.App, uc UsecasesForCreateChat) {
	// Тело запроса для создания чата.
	type requestBody struct {
		Name string `json:"name"`
	}
	router.Post(
		"/chats",
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := createChat.In{
				ChiefUserID: Session(context).UserID,
				Name:        rb.Name,
			}

			out, err := uc.CreateChat(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(uc),
	)
}

type UsecasesForCreateChat interface {
	CreateChat(createChat.In) (createChat.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
