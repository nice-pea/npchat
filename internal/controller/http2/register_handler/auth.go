package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	"github.com/nice-pea/npchat/internal/service"
)

// LoginByPassword регистрирует обработчик, позволяющий авторизоваться по логину и паролю.
// Доступен без предварительной аутентификации (публичная цепочка middleware).
//
// Метод: POST /auth/password/login
func LoginByPassword(router *fiber.App, ss Services) {
	// Тело запроса для авторизации по логину и паролю.
	type requestBody struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	router.Post(
		"/auth/password/login",
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := service.BasicAuthLoginIn{
				Login:    rb.Login,
				Password: rb.Password,
			}

			out, err := ss.Users().BasicAuthLogin(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
	)
}

// RegistrationByPassword регистрирует обработчик, позволяющий регистрироваться по логину и паролю.
// Доступен без предварительной аутентификации (публичная цепочка middleware).
//
// Метод: POST /auth/password/registration
func RegistrationByPassword(router *fiber.App, ss Services) {
	// Тело запроса для авторизации по логину и паролю.
	type requestBody struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Nick     string `json:"nick"`
	}
	router.Post(
		"/auth/password/registration",
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := service.BasicAuthRegistrationIn{
				Login:    rb.Login,
				Password: rb.Password,
				Name:     rb.Name,
				Nick:     rb.Nick,
			}

			out, err := ss.Users().BasicAuthRegistration(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
		recover2.New(),
	)
}
