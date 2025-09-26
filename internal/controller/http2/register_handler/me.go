package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	userProfile "github.com/nice-pea/npchat/internal/usecases/users/user_profile"
)

// Me регистрирует HTTP-обработчик для получения информации о пользователе по токену.
// Данный обработчик доступен только авторизованным пользователям.
//
// Метод: GET /chats
func Me(router *fiber.App, uc UsecasesForGetUser, jwtParser middleware.JwtParser) {
	router.Get(
		"/me",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
		func(ctx *fiber.Ctx) error {
			input := userProfile.In{
				SubjectID: UserID(ctx),
				UserID:    UserID(ctx),
			}

			out, err := uc.UserProfile(input)
			if err != nil {
				return err
			}

			return ctx.JSON(out)
		},
	)
}

// UsecasesForMe определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForMe interface {
	UserProfile(userProfile.In) (userProfile.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
