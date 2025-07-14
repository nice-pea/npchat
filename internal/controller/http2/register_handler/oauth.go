package register_handler

import (
	"errors"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nice-pea/npchat/internal/service"
)

// OAuthInitRegistration регистрирует обработчик, инициирующий процесс регистрации через OAuth.
// Редиректит пользователя на страницу авторизации провайдера.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/registration
func OAuthInitRegistration(router *fiber.App, ss services) {
	router.Get(
		"/oauth/:provider/registration",
		func(context *fiber.Ctx) error {
			// Формируем входные данные для инициализации OAuth-регистрации
			input := service.InitOAuthRegistrationIn{
				Provider: context.Params("provider"), // Получаем имя провайдера из URL
			}

			// Инициируем OAuth-процесс регистрации у сервиса
			out, err := ss.Users().InitOAuthRegistration(input)
			if err != nil {
				return err
			}

			// Сохраняем параметр state из URL в куке для последующей проверки безопасности
			if err = setOAuthCookie(context, out.RedirectURL); err != nil {
				return err
			}

			// Возвращаем команду редиректа на сторону провайдера
			return context.Redirect(out.RedirectURL, fiber.StatusTemporaryRedirect)
		},
	)
}

// OAuthCompleteRegistrationCallback регистрирует обработчик, завершающий регистрацию через OAuth.
// Обрабатывает callback от провайдера после успешной авторизации.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/registration/callback
func OAuthCompleteRegistrationCallback(router *fiber.App, ss services) {
	router.Get(
		"/oauth/:provider/registration/callback",
		func(context *fiber.Ctx) error {
			// Проверяем, что запрос пришёл из доверенного источника через сравнение state
			if err := validateOAuthCookie(context); err != nil {
				return err
			}

			input := service.CompeteOAuthRegistrationIn{
				UserCode: context.Query("code"),
				Provider: context.Params("provider"),
			}

			out, err := ss.Users().CompeteOAuthRegistration(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// OAuthInitLogin регистрирует обработчик, инициирующий процесс входа через OAuth.
// Редиректит пользователя на страницу авторизации провайдера.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/login
func OAuthInitLogin(router *fiber.App, ss services) {
	router.Get(
		"/oauth/:provider/login",
		func(context *fiber.Ctx) error {
			input := service.InitOAuthLoginIn{
				Provider: context.Params("provider"), // Получаем имя провайдера из URL
			}

			// Инициируем OAuth-процесс входа у сервиса
			out, err := ss.Users().InitOAuthLogin(input)
			if err != nil {
				return err
			}

			// Сохраняем параметр state из URL в куке для последующей проверки безопасности
			if err = setOAuthCookie(context, out.RedirectURL); err != nil {
				return err
			}

			// Возвращаем команду редиректа на сторону провайдера
			return context.Redirect(out.RedirectURL, fiber.StatusTemporaryRedirect)
		},
	)
}

// OAuthCompleteLoginCallback регистрирует обработчик, завершающий процесс входа через OAuth.
// Обрабатывает callback от провайдера после успешной авторизации.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/login/callback
func OAuthCompleteLoginCallback(router *fiber.App, ss services) {
	router.Get(
		"/oauth/:provider/login/callback",
		func(context *fiber.Ctx) error {
			// Проверяем, что запрос пришёл из доверенного источника через сравнение state
			if err := validateOAuthCookie(context); err != nil {
				return err
			}

			input := service.CompleteOAuthLoginIn{
				UserCode: context.Query("code"),
				Provider: context.Params("provider"),
			}

			out, err := ss.Users().CompleteOAuthLogin(input)
			if err != nil {
				return err
			}

			return context.JSON(out)
		},
	)
}

// oauthCookieName — имя куки, в которую сохраняется параметр state для защиты от CSRF
const oauthCookieName = "oauthState"

// setOAuthCookie устанавливает куку с параметром state из строки редиректа
func setOAuthCookie(context *fiber.Ctx, redirectURL string) error {
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

// validateOAuthCookie проверяет, что значение 'state' в запросе совпадает с тем, что было сохранено в куке.
// Это защищает от CSRF-атак.
func validateOAuthCookie(context *fiber.Ctx) error {
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
