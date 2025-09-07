package oauthProvider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/nice-pea/npchat/internal/domain/userr"
)

// Google представляет собой структуру для работы с Oauth аутентификацией через Google.
type Google struct {
	config *oauth2.Config // Конфигурация Oauth для Google
}

type GoogleConfig struct {
	ClientID     string // Идентификатор клиента для Oauth
	ClientSecret string // Секрет клиента для Oauth
	RedirectURL  string // URL для перенаправления после аутентификации
}

func NewGoogle(cfg GoogleConfig) *Google {
	return &Google{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,     // Идентификатор клиента
			ClientSecret: cfg.ClientSecret, // Секрет клиента
			Endpoint:     google.Endpoint,  // Использует конечную точку Google для Oauth
			RedirectURL:  cfg.RedirectURL,  // URL для перенаправления
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",   // Запрашивает доступ к электронной почте пользователя
				"https://www.googleapis.com/auth/userinfo.profile", // Запрашивает доступ к профилю пользователя
			},
		},
	}
}

// Name возвращает имя провайдера Oauth.
func (o *Google) Name() string {
	return "google"
}

// Exchange обменивает код авторизации на токен Oauth.
func (o *Google) Exchange(code string) (userr.OpenAuthToken, error) {
	// Обменять код авторизации на токен Oauth
	token, err := o.config.Exchange(context.Background(), code)
	if err != nil {
		return userr.OpenAuthToken{}, err
	}

	return userr.NewOpenAuthToken(
		token.AccessToken,  // Токен доступа
		token.Type(),       // Тип токена
		token.RefreshToken, // Токен обновления
		token.Expiry,       // Время истечения токена
	)
}

// User получает информацию о пользователе Google, используя токен Oauth.
func (o *Google) User(token userr.OpenAuthToken) (userr.OpenAuthUser, error) {
	// Выполить GET-запрос для получения информации о пользователе
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return userr.OpenAuthUser{}, err // Возвращает ошибку, если запрос не удался
	}

	// Прочитать данные ответ
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return userr.OpenAuthUser{}, err // Возвращает ошибку, если чтение не удалось
	}
	_ = response.Body.Close()

	// Сложить ответ в структуру
	var googleUser struct {
		ID      string `json:"id"`      // Идентификатор пользователя
		Email   string `json:"email"`   // Электронная почта пользователя
		Name    string `json:"name"`    // Имя пользователя
		Picture string `json:"picture"` // URL изображения профиля пользователя
	}
	if err = json.Unmarshal(data, &googleUser); err != nil {
		return userr.OpenAuthUser{}, err // Возвращает ошибку, если разбор JSON не удался
	}

	return userr.NewOpenAuthUser(
		googleUser.ID,      // Идентификатор пользователя
		o.Name(),           // Имя провайдера
		googleUser.Email,   // Электронная почта пользователя
		googleUser.Name,    // Имя пользователя
		googleUser.Picture, // URL изображения профиля
		token,
	)
}

// AuthorizationURL генерирует URL для авторизации с использованием кода состояния.
func (o *Google) AuthorizationURL(state string) string {
	// Сгенерировать URL
	return o.config.AuthCodeURL(state)
}
