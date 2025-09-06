package register_handler

import (
	"errors"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	oauthAuthorize "github.com/nice-pea/npchat/internal/usecases/users/oauth/oauth_authorize"
	oauthComplete "github.com/nice-pea/npchat/internal/usecases/users/oauth/oauth_complete"
)

// OauthAuthorize регистрирует обработчик, инициирующий процесс регистрации через Oauth.
// Редиректит пользователя на страницу авторизации провайдера.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/authorize
func OauthAuthorize(router *fiber.App, uc UsecasesForOauthAuthorize) {
	router.Get(
		"/oauth/:provider/authorize",
		recover2.New(),
		func(ctx *fiber.Ctx) error {
			// Проверяем, что запрос пришёл из доверенного источника через сравнение state
			if err := validateOauthCookie(ctx); err != nil {
				return err
			}

			// Формируем входные данные для инициализации Oauth-регистрации
			input := oauthAuthorize.In{
				Provider:         ctx.Params("provider"), // Получаем имя провайдера из URL
				CompleteCallback: ctx.Get("Origin") + "/oauth/" + ctx.Params("provider") + "/callback",
			}

			// Инициируем Oauth-процесс регистрации у сервиса
			out, err := uc.OauthAuthorize(input)
			if err != nil {
				return err
			}

			// Сохраняем параметр state из URL в куке для последующей проверки безопасности
			if err = setOauthCookie(ctx, out.RedirectURL); err != nil {
				return err
			}

			// Возвращаем команду редиректа на сторону провайдера
			return ctx.Redirect(out.RedirectURL, fiber.StatusTemporaryRedirect)
		},
	)
}

// UsecasesForOauthAuthorize определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForOauthAuthorize interface {
	OauthAuthorize(oauthAuthorize.In) (oauthAuthorize.Out, error)
}

// OauthCallback регистрирует обработчик, завершающий регистрацию через Oauth.
// Обрабатывает callback от провайдера после успешной авторизации.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/callback
func OauthCallback(router *fiber.App, uc UsecasesForOauthCallback) {
	router.Get(
		"/oauth/:provider/callback",
		recover2.New(),
		func(ctx *fiber.Ctx) error {
			// Проверяем, что запрос пришёл из доверенного источника через сравнение state
			if err := validateOauthCookie(ctx); err != nil {
				return err
			}

			input := oauthComplete.In{
				UserCode: ctx.Query("code"),
				Provider: ctx.Params("provider"),
			}

			out, err := uc.OauthComplete(input)
			if err != nil {
				return err
			}

			return ctx.JSON(out)
		},
	)
}

// UsecasesForOauthCallback определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForOauthCallback interface {
	OauthComplete(oauthComplete.In) (oauthComplete.Out, error)
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
