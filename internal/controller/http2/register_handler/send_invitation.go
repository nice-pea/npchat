package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	"github.com/nice-pea/npchat/internal/service"
)

// SendInvitation регистрирует обработчик, позволяющий отправить приглашение в чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /invitations
func SendInvitation(router *fiber.App, ss Services) {
	// Тело запроса для отправки приглашения.
	type requestBody struct {
		ChatID uuid.UUID `json:"chat_id"`
		UserID uuid.UUID `json:"user_id"`
	}
	router.Post(
		"/invitations",
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := service.SendInvitationIn{
				SubjectID: Session(context).UserID,
				ChatID:    rb.ChatID,
				UserID:    rb.UserID,
			}

			out, err := ss.Chats().SendInvitation(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
		middleware.RequireAuthorizedSession(ss.Sessions()),
	)
}
