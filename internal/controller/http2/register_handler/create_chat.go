package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// CreateChat регистрирует обработчик, позволяющий создать новый чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /chats
func CreateChat(router *fiber.App, ss Services) {
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

			input := service.CreateChatIn{
				ChiefUserID: Session(context).UserID,
				Name:        rb.Name,
			}

			out, err := ss.Chats().CreateChat(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(ss.Sessions()),
	)
}
