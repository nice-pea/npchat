package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	updateName "github.com/nice-pea/npchat/internal/service/chats/update_name"
)

// UpdateChatName регистрирует обработчик, позволяющий обновить название чата.
// Доступен только авторизованным пользователям, которые являются главными администраторами чата.
//
// Метод: PUT /chats/{chatID}/name
func UpdateChatName(router *fiber.App, uc UsecasesForUpdateName) {
	// Тело запроса для обновления названия чата.
	type requestBody struct {
		NewName string `json:"new_name"`
	}
	router.Put(
		"/chats/:chatID/name",
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := updateName.In{
				SubjectID: Session(context).UserID,
				ChatID:    ParamsUUID(context, "chatID"),
				NewName:   rb.NewName,
			}

			out, err := uc.UpdateName(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(uc),
	)
}

// UsecasesForUpdateName определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForUpdateName interface {
	UpdateName(updateName.In) (updateName.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
