package register_handler

import (
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
	recover2 "github.com/gofiber/fiber/v2/middleware/recover"

	oauthAuthorize "github.com/nice-pea/npchat/internal/usecases/users/oauth/oauth_authorize"
	oauthComplete "github.com/nice-pea/npchat/internal/usecases/users/oauth/oauth_complete"
)

// OauthAuthorize регистрирует обработчик, инициирующий процесс входа через Oauth.
// Перенаправляет пользователя на страницу авторизации провайдера.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/authorize
func OauthAuthorize(router *fiber.App, uc UsecasesForOauthAuthorize) {
	router.Get(
		"/oauth/:provider/authorize",
		recover2.New(),
		func(ctx *fiber.Ctx) error {
			// Формируем входные данные для инициализации Oauth-регистрации
			input := oauthAuthorize.In{
				Provider: ctx.Params("provider"), // Получаем имя провайдера из URL
			}

			// Инициируем Oauth-процесс регистрации у сервиса
			out, err := uc.OauthAuthorize(input)
			if err != nil {
				return err
			}

			// Сохраняем state в куке для последующей проверки безопасности
			setOauthCookie(ctx, out.State)

			// Возвращаем команду редиректа на сторону провайдера
			return ctx.Redirect(out.RedirectURL, fiber.StatusTemporaryRedirect)
		},
	)
}

// UsecasesForOauthAuthorize определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForOauthAuthorize interface {
	OauthAuthorize(oauthAuthorize.In) (oauthAuthorize.Out, error)
}

// OauthCallback регистрирует обработчик, завершающий вход через Oauth.
// Обрабатывает callback от провайдера после успешной авторизации.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/callback
func OauthCallback(router *fiber.App, uc UsecasesForOauthCallback, jwtIssuer JwtIssuer) {
	router.Get(
		"/oauth/:provider/callback",
		recover2.New(),
		func(ctx *fiber.Ctx) error {
			// Проверяем, что запрос пришёл из доверенного источника через сравнение state
			if err := validateOauthCookie(ctx); err != nil {
				return err
			}
			// Удалить куку
			clearOauthCookie(ctx)

			input := oauthComplete.In{
				UserCode: ctx.Query("code"),
				Provider: ctx.Params("provider"),
			}

			out, err := uc.OauthComplete(input)
			if err != nil {
				return err
			}

			token, err := jwtIssuer.Issue(out.Session)
			if err != nil {
				return err
			}

			return ctx.JSON(fiber.Map{
				"Out": out,
				"Jwt": token,
			})
		},
	)
}

// UsecasesForOauthCallback определяет интерфейс для доступа к сценариям использования бизнес-логики
type UsecasesForOauthCallback interface {
	OauthComplete(oauthComplete.In) (oauthComplete.Out, error)
}

// oauthCookieName возвращает имя куки, в которую сохраняется параметр state для защиты от CSRF
func oauthCookieName(ctx *fiber.Ctx) string {
	return ctx.Params("provider") + "-oauthState"
}

// clearOauthCookie очищает одноразовую state-куку
func clearOauthCookie(ctx *fiber.Ctx) {
	ctx.ClearCookie(oauthCookieName(ctx))
}

// setOauthCookie устанавливает куку с параметром state
func setOauthCookie(ctx *fiber.Ctx, state string) {
	// Устанавливаем куку с этим значением
	ctx.Cookie(&fiber.Cookie{
		Name:     oauthCookieName(ctx),
		Value:    state,
		Expires:  time.Now().Add(time.Minute * 15), // Время жизни кука
		HTTPOnly: true,                             // Защита от XSS
		Secure:   true,                             // Только по HTTPS
		Path:     "/",                              // Доступна по всему домену
	})
}

// errWrongState — ошибка, возникающая при несовпадении значения state в куке и запросе
var errWrongState = errors.New("неправильный state")

// validateOauthCookie проверяет, что значение 'state' в запросе совпадает с тем, что было сохранено в куке.
// Это защищает от CSRF-атак.
func validateOauthCookie(ctx *fiber.Ctx) error {
	// Получаем куку с именем oauthState
	oauthState := ctx.Cookies(oauthCookieName(ctx))

	if oauthState == "" {
		// Если куки нет — возвращаем ошибку несоответствия state
		return errWrongState
	}

	// Сравниваем значение state из запроса с тем, что храним в куке
	if ctx.Query("state") != oauthState {
		return errWrongState
	}

	return nil
}
