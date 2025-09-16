package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	sendInvitation "github.com/nice-pea/npchat/internal/usecases/chats/send_invitation"
)

// SendInvitation регистрирует обработчик, позволяющий отправить приглашение в чат.
// Доступен только авторизованным пользователям.
//
// Метод: POST /invitations
func SendInvitation(router *fiber.App, uc UsecasesForSendInvitation, jparser middleware.JwtParser) {
	// Тело запроса для отправки приглашения.
	type requestBody struct {
		ChatID uuid.UUID `json:"chat_id"`
		UserID uuid.UUID `json:"user_id"`
	}
	router.Post(
		"/invitations",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jparser),
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := sendInvitation.In{
				SubjectID: UserID(context),
				ChatID:    rb.ChatID,
				UserID:    rb.UserID,
			}

			out, err := uc.SendInvitation(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// UsecasesForSendInvitation определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForSendInvitation interface {
	SendInvitation(sendInvitation.In) (sendInvitation.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
