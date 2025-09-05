package register_handler

import (
	"errors"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	completeOauthLogin "github.com/nice-pea/npchat/internal/usecases/users/oauth/complete_oauth_login"
	completeOauthRegistration "github.com/nice-pea/npchat/internal/usecases/users/oauth/complete_oauth_registration"
	initOauthLogin "github.com/nice-pea/npchat/internal/usecases/users/oauth/init_oauth_login"
	initOauthRegistration "github.com/nice-pea/npchat/internal/usecases/users/oauth/init_oauth_registration"
)

// OauthInitRegistration регистрирует обработчик, инициирующий процесс регистрации через Oauth.
// Редиректит пользователя на страницу авторизации провайдера.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/registration
func OauthInitRegistration(router *fiber.App, uc UsecasesForOauthInitRegistration) {
	router.Get(
		"/oauth/:provider/registration",
		recover2.New(),
		func(context *fiber.Ctx) error {
			// Формируем входные данные для инициализации Oauth-регистрации
			input := initOauthRegistration.In{
				Provider: context.Params("provider"), // Получаем имя провайдера из URL
			}

			// Инициируем Oauth-процесс регистрации у сервиса
			out, err := uc.InitOauthRegistration(input)
			if err != nil {
				return err
			}

			// Сохраняем параметр state из URL в куке для последующей проверки безопасности
			if err = setOauthCookie(context, out.RedirectURL); err != nil {
				return err
			}

			// Возвращаем команду редиректа на сторону провайдера
			return context.Redirect(out.RedirectURL, fiber.StatusTemporaryRedirect)
		},
	)
}

// UsecasesForOauthInitRegistration определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForOauthInitRegistration interface {
	InitOauthRegistration(initOauthRegistration.In) (initOauthRegistration.Out, error)
}

// OauthCompleteRegistrationCallback регистрирует обработчик, завершающий регистрацию через Oauth.
// Обрабатывает callback от провайдера после успешной авторизации.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/registration/callback
func OauthCompleteRegistrationCallback(router *fiber.App, uc UsecasesForOauthCompleteRegistrationCallback) {
	router.Get(
		"/oauth/:provider/registration/callback",
		recover2.New(),
		func(context *fiber.Ctx) error {
			// Проверяем, что запрос пришёл из доверенного источника через сравнение state
			if err := validateOauthCookie(context); err != nil {
				return err
			}

			input := completeOauthRegistration.In{
				UserCode: context.Query("code"),
				Provider: context.Params("provider"),
			}

			out, err := uc.CompleteOauthRegistration(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// UsecasesForOauthCompleteRegistrationCallback определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForOauthCompleteRegistrationCallback interface {
	CompleteOauthRegistration(completeOauthRegistration.In) (completeOauthRegistration.Out, error)
}

// OauthInitLogin регистрирует обработчик, инициирующий процесс входа через Oauth.
// Редиректит пользователя на страницу авторизации провайдера.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/login
func OauthInitLogin(router *fiber.App, uc UsecasesForOauthInitLogin) {
	router.Get(
		"/oauth/:provider/login",
		recover2.New(),
		func(context *fiber.Ctx) error {
			input := initOauthLogin.In{
				Provider: context.Params("provider"), // Получаем имя провайдера из URL
			}

			// Инициируем Oauth-процесс входа у сервиса
			out, err := uc.InitOauthLogin(input)
			if err != nil {
				return err
			}

			// Сохраняем параметр state из URL в куке для последующей проверки безопасности
			if err = setOauthCookie(context, out.RedirectURL); err != nil {
				return err
			}

			// Возвращаем команду редиректа на сторону провайдера
			return context.Redirect(out.RedirectURL, fiber.StatusTemporaryRedirect)
		},
	)
}

// UsecasesForOauthInitLogin определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForOauthInitLogin interface {
	InitOauthLogin(initOauthLogin.In) (initOauthLogin.Out, error)
}

// OauthCompleteLoginCallback регистрирует обработчик, завершающий процесс входа через Oauth.
// Обрабатывает callback от провайдера после успешной авторизации.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/login/callback
func OauthCompleteLoginCallback(router *fiber.App, uc UsecasesForOauthCompleteLoginCallback) {
	router.Get(
		"/oauth/:provider/login/callback",
		recover2.New(),
		func(context *fiber.Ctx) error {
			// Проверяем, что запрос пришёл из доверенного источника через сравнение state
			if err := validateOauthCookie(context); err != nil {
				return err
			}

			input := completeOauthLogin.In{
				UserCode: context.Query("code"),
				Provider: context.Params("provider"),
			}

			out, err := uc.CompleteOauthLogin(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// UsecasesForOauthCompleteLoginCallback определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForOauthCompleteLoginCallback interface {
	CompleteOauthLogin(completeOauthLogin.In) (completeOauthLogin.Out, error)
}

// oauthCookieName — имя куки, в которую сохраняется параметр state для защиты от CSRF
const oauthCookieName = "oauthState"

// setOauthCookie устанавливает куку с параметром state из строки редиректа
func setOauthCookie(context *fiber.Ctx, redirectURL string) error {
	// Парсим URL, чтобы получить query-параметры
	parsedUrl, err := url.Parse(redirectURL)
	if err != nil {
		return err
	}

	// Получаем значение state из URL
	state := parsedUrl.Query().Get("state")

	// Устанавливаем куку с этим значением
	context.Cookie(&fiber.Cookie{
		Name:     oauthCookieName,
		Value:    state,
		Expires:  time.Now().Add(time.Hour), // Кука живёт 1 час
		HTTPOnly: true,                      // Защита от XSS
		Secure:   true,                      // Только по HTTPS
		Path:     "/",                       // Доступна по всему домену
	})

	return nil
}

// errWrongState — ошибка, возникающая при несовпадении значения state в куке и запросе
var errWrongState = errors.New("неправильный state")

// validateOauthCookie проверяет, что значение 'state' в запросе совпадает с тем, что было сохранено в куке.
// Это защищает от CSRF-атак.
func validateOauthCookie(context *fiber.Ctx) error {
	// Получаем куку с именем oauthState
	oauthState := context.Cookies(oauthCookieName)

	if oauthState == "" {
		// Если куки нет — возвращаем ошибку несоответствия state
		return errWrongState
	}

	// Сравниваем значение state из запроса с тем, что храним в куке
	if context.Query("state") != oauthState {
		return errWrongState
	}

	return nil
}
