package oauth_provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

// Google представляет собой структуру для работы с OAuth2 аутентификацией через Google.
type Google struct {
	ClientID     string // Идентификатор клиента для OAuth2
	ClientSecret string // Секрет клиента для OAuth2
	RedirectURL  string // URL для перенаправления после аутентификации
}

// Name возвращает имя провайдера OAuth.
func (o *Google) Name() string {
	return "google"
}

// config создает и возвращает конфигурацию OAuth2 для Google.
func (o *Google) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.ClientID,      // Идентификатор клиента
		ClientSecret: o.ClientSecret,  // Секрет клиента
		Endpoint:     google.Endpoint, // Использует конечную точку Google для OAuth2
		RedirectURL:  o.RedirectURL,   // URL для перенаправления
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",   // Запрашивает доступ к электронной почте пользователя
			"https://www.googleapis.com/auth/userinfo.profile", // Запрашивает доступ к профилю пользователя
		},
	}
}

// Exchange обменивает код авторизации на токен OAuth.
func (o *Google) Exchange(code string) (userr.OpenAuthToken, error) {
	// Обменять код авторизации на токен OAuth
	token, err := o.config().Exchange(context.Background(), code)
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

// User получает информацию о пользователе Google, используя токен OAuth.
func (o *Google) User(token userr.OpenAuthToken) (userr.OpenAuthUser, error) {
	const getUser = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	// Выполняет GET-запрос для получения информации о пользователе
	response, err := http.Get(getUser + token.AccessToken)
	if err != nil {
		return userr.OpenAuthUser{}, err // Возвращает ошибку, если запрос не удался
	}

	defer func() { _ = response.Body.Close() }() // Закрывает тело ответа после завершения работы с ним
	data, err := io.ReadAll(response.Body)       // Читает данные из тела ответа
	if err != nil {
		return userr.OpenAuthUser{}, err // Возвращает ошибку, если чтение не удалось
	}

	// Сложить данные в структуру ответа
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
	return o.config().AuthCodeURL(state) // Генерирует URL для авторизации
}
