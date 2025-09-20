package register_handler

import (
	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	basicAuthLogin "github.com/nice-pea/npchat/internal/usecases/users/basic_auth/basic_auth_login"
	basicAuthRegistration "github.com/nice-pea/npchat/internal/usecases/users/basic_auth/basic_auth_registration"
)

// LoginByPassword регистрирует обработчик, позволяющий авторизоваться по логину и паролю.
// Доступен без предварительной аутентификации (публичная цепочка middleware).
//
// Метод: POST /auth/password/login
func LoginByPassword(router *fiber.App, uc UsecasesForLoginByPassword, jwtIssuer JwtIssuer) {
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

			input := basicAuthLogin.In{
				Login:    rb.Login,
				Password: rb.Password,
			}

			out, err := uc.BasicAuthLogin(input)
			if err != nil {
				return err
			}

			token, err := jwtIssuer.Issue(out.Session)

			if err != nil {
				// TODO: тут мб если сессии без ошибок создались то не возвращать ошибку
				return err
			}

			return context.JSON(LoginByPasswordOut{out, token})
		},
		recover2.New(),
	)
}

type LoginByPasswordOut struct {
	Out basicAuthLogin.Out
	Jwt string
}

// UsecasesForAcceptInvitation определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForLoginByPassword interface {
	BasicAuthLogin(basicAuthLogin.In) (basicAuthLogin.Out, error)
}

// RegistrationByPassword регистрирует обработчик, позволяющий регистрироваться по логину и паролю.
// Доступен без предварительной аутентификации (публичная цепочка middleware).
//
// Метод: POST /auth/password/registration
func RegistrationByPassword(router *fiber.App, uc UsecasesForRegistrationByPassword) {
	// Тело запроса для авторизации по логину и паролю.
	type requestBody struct {
		Login    string `json:"login"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Nick     string `json:"nick"`
	}
	router.Post(
		"/auth/password/registration",
		recover2.New(),
		func(context *fiber.Ctx) error {
			var rb requestBody
			// Декодируем тело запроса в структуру requestBody.
			if err := context.BodyParser(&rb); err != nil {
				return err
			}

			input := basicAuthRegistration.In{
				Login:    rb.Login,
				Password: rb.Password,
				Name:     rb.Name,
				Nick:     rb.Nick,
			}

			out, err := uc.BasicAuthRegistration(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// UsecasesForRegistrationByPassword определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForRegistrationByPassword interface {
	BasicAuthRegistration(basicAuthRegistration.In) (basicAuthRegistration.Out, error)
}
