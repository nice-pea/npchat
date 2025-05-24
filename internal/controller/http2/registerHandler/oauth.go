package registerHandler

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// OAuthInitRegistration регистрирует обработчик, инициирующий процесс регистрации через OAuth.
// Редиректит пользователя на страницу авторизации провайдера.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/registration
func OAuthInitRegistration(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/registration",
		middleware.EmptyChain, // Нет аутентификации на этом этапе
		func(context http2.Context) (any, error) {
			// Формируем входные данные для инициализации OAuth-регистрации
			input := service.OAuthInitRegistrationInput{
				Provider: http2.PathStr(context, "provider"), // Получаем имя провайдера из URL
			}

			// Инициируем OAuth-процесс регистрации у сервиса
			out, err := context.Services().OAuth().InitRegistration(input)
			if err != nil {
				return nil, err
			}

			// Сохраняем параметр state из URL в куке для последующей проверки безопасности
			if err = setOAuthCookie(context, out.RedirectURL); err != nil {
				return nil, err
			}

			// Возвращаем команду редиректа на сторону провайдера
			return http2.Redirect{
				URL:  out.RedirectURL,
				Code: http.StatusTemporaryRedirect,
			}, nil
		},
	)
}

// OAuthCompleteRegistrationCallback регистрирует обработчик, завершающий регистрацию через OAuth.
// Обрабатывает callback от провайдера после успешной авторизации.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/registration/callback
func OAuthCompleteRegistrationCallback(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/registration/callback",
		middleware.EmptyChain, // Нет аутентификации — пользователь ещё не зарегистрирован
		func(context http2.Context) (any, error) {
			// Проверяем, что запрос пришёл из доверенного источника через сравнение state
			if err := validateOAuthCookie(context); err != nil {
				return nil, err
			}

			input := service.OAuthCompeteRegistrationInput{
				UserCode: http2.FormStr(context, "code"),
				Provider: http2.PathStr(context, "provider"),
			}

			return context.Services().OAuth().CompeteRegistration(input)
		},
	)
}

// OAuthInitLogin регистрирует обработчик, инициирующий процесс входа через OAuth.
// Редиректит пользователя на страницу авторизации провайдера.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/login
func OAuthInitLogin(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/login",
		middleware.EmptyChain, // Нет аутентификации на этом этапе
		func(context http2.Context) (any, error) {
			// Формируем входные данные для инициализации OAuth-входа
			input := service.OAuthInitLoginInput{
				Provider: http2.PathStr(context, "provider"), // Получаем имя провайдера из URL
			}

			// Инициируем OAuth-процесс входа у сервиса
			out, err := context.Services().OAuth().InitLogin(input)
			if err != nil {
				return nil, err
			}

			// Сохраняем параметр state из URL в куке для последующей проверки безопасности
			if err = setOAuthCookie(context, out.RedirectURL); err != nil {
				return nil, err
			}

			// Возвращаем команду редиректа на сторону провайдера
			return http2.Redirect{
				URL:  out.RedirectURL,
				Code: http.StatusTemporaryRedirect,
			}, nil
		},
	)
}

// OAuthCompleteLoginCallback регистрирует обработчик, завершающий процесс входа через OAuth.
// Обрабатывает callback от провайдера после успешной авторизации.
// Данный обработчик не требует аутентификации.
//
// Метод: GET /oauth/{provider}/login/callback
func OAuthCompleteLoginCallback(router http2.Router) {
	router.HandleFunc(
		"GET /oauth/{provider}/login/callback",
		middleware.EmptyChain, // Нет аутентификации — пользователь ещё не вошёл
		func(context http2.Context) (any, error) {
			// Проверяем, что запрос пришёл из доверенного источника через сравнение state
			if err := validateOAuthCookie(context); err != nil {
				return nil, err
			}

			// Формируем входные данные для завершения OAuth-входа
			input := service.OAuthCompleteLoginInput{
				UserCode: http2.FormStr(context, "code"),     // Код, переданный провайдером
				Provider: http2.PathStr(context, "provider"), // Имя провайдера из URL
			}

			// Завершаем вход через OAuth-сервис
			return context.Services().OAuth().CompleteLogin(input)
		},
	)
}

// oauthCookieName — имя куки, в которую сохраняется параметр state для защиты от CSRF
const oauthCookieName = "oauthState"

// setOAuthCookie устанавливает куку с параметром state из строки редиректа
func setOAuthCookie(context http2.Context, redirectURL string) error {
	// Парсим URL, чтобы получить query-параметры
	parsedUrl, err := url.Parse(redirectURL)
	if err != nil {
		return err
	}

	// Получаем значение state из URL
	state := parsedUrl.Query().Get("state")

	// Устанавливаем куку с этим значением
	http.SetCookie(context.Writer(), &http.Cookie{
		Name:     oauthCookieName,
		Value:    state,
		Expires:  time.Now().Add(time.Hour), // Кука живёт 1 час
		HttpOnly: true,                      // Защита от XSS
		Secure:   true,                      // Только по HTTPS
		Path:     "/",                       // Доступна по всему домену
	})

	return nil
}

// errWrongState — ошибка, возникающая при несовпадении значения state в куке и запросе
var errWrongState = errors.New("неправильный state")

// validateOAuthCookie проверяет, что значение 'state' в запросе совпадает с тем, что было сохранено в куке.
// Это защищает от CSRF-атак.
func validateOAuthCookie(context http2.Context) error {
	// Получаем куку с именем oauthState
	oauthState, err := context.Request().Cookie(oauthCookieName)

	if errors.Is(err, http.ErrNoCookie) {
		// Если куки нет — возвращаем ошибку несоответствия state
		return errWrongState
	} else if err != nil {
		// Любые другие ошибки также возвращаем
		return err
	}

	// Сравниваем значение state из запроса с тем, что храним в куке
	if http2.FormStr(context, "state") != oauthState.Value {
		return errWrongState
	}

	return nil
}
