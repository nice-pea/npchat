package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	updateName "github.com/nice-pea/npchat/internal/usecases/chats/update_name"
)

// UpdateChatName регистрирует обработчик, позволяющий обновить название чата.
// Доступен только авторизованным пользователям, которые являются главными администраторами чата.
//
// Метод: PUT /chats/{chatID}/name
func UpdateChatName(router *fiber.App, uc UsecasesForUpdateName, jwtParser middleware.JwtParser) {
	// Тело запроса для обновления названия чата.
	type requestBody struct {
		NewName string `json:"new_name"`
	}
	router.Put(
		"/chats/:chatID/name",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
		func(ctx *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := ctx.BodyParser(&rb); err != nil {
				return err
			}

			input := updateName.In{
				SubjectID: UserID(ctx),
				ChatID:    ParamsUUID(ctx, "chatID"),
				NewName:   rb.NewName,
			}

			out, err := uc.UpdateName(input)
			if err != nil {
				return err
			}

			return ctx.JSON(out)
		},
	)
}

// UsecasesForUpdateName определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForUpdateName interface {
	UpdateName(updateName.In) (updateName.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
