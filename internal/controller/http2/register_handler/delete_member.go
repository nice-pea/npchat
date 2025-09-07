package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	deleteMember "github.com/nice-pea/npchat/internal/usecases/chats/delete_member"
)

// DeleteMember регистрирует обработчик, позволяющий удалить участника из чата.
// Доступен только авторизованным пользователям, которые являются главными администраторами чата.
//
// Метод: DELETE /chats/{chatID}/members
func DeleteMember(router *fiber.App, uc UsecasesForDeleteMember) {
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

			input := deleteMember.In{
				SubjectID: Session(context).UserID,
				ChatID:    ParamsUUID(context, "chatID"),
				UserID:    rb.UserID,
			}

			out, err := uc.DeleteMember(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, nil),
	)
}

// UsecasesForDeleteMember определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForDeleteMember interface {
	DeleteMember(deleteMember.In) (deleteMember.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
