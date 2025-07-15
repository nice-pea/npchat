package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// DeleteMember регистрирует обработчик, позволяющий удалить участника из чата.
// Доступен только авторизованным пользователям, которые являются главными администраторами чата.
//
// Метод: DELETE /chats/{chatID}/members
func DeleteMember(router *fiber.App, ss Services) {
	// Тело запроса для удаления участника из чата.
	type requestBody struct {
		UserID uuid.UUID `json:"user_id"`
	}
	router.Delete(
		"/chats/:chatID/members",
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := service.DeleteMemberIn{
				SubjectID: Session(context).UserID,
				ChatID:    ParamsUUID(context, "chatID"),
				UserID:    rb.UserID,
			}

			return ss.Chats().DeleteMember(input)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(ss.Sessions()),
	)
}
