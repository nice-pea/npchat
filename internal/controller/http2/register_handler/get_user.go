package registerHandler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/controller/http2/middleware"
	userProfile "github.com/nice-pea/npchat/internal/usecases/users/user_profile"
)

// GetUser регистрирует HTTP-обработчик для получения информации о пользователе.
// Данный обработчик доступен только авторизованным пользователям.
//
// Метод: GET /chats
func GetUser(router *fiber.App, uc UsecasesForGetUser, jwtParser middleware.JwtParser) {
	router.Get(
		"/users/{id}",
		recover2.New(),
		middleware.RequireAuthorizedSession(uc, jwtParser),
		func(ctx *fiber.Ctx) error {
			input := userProfile.In{
				SubjectID: UserID(ctx),
				UserID:    ParamsUUID(ctx, "id"),
			}

			out, err := uc.UserProfile(input)
			if err != nil {
				return err
			}

			return ctx.JSON(out)
		},
	)
}

// UsecasesForGetUser определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForGetUser interface {
	UserProfile(userProfile.In) (userProfile.Out, error)
	middleware.UsecasesForRequireAuthorizedSession
}
