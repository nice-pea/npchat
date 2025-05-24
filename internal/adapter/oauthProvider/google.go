package oauthProvider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/saime-0/nice-pea-chat/internal/domain"
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
func (o *Google) Exchange(code string) (domain.OAuthToken, error) {
	// Обменять код авторизации на токен OAuth
	token, err := o.config().Exchange(context.Background(), code)
	if err != nil {
		return domain.OAuthToken{}, err
	}

	return domain.OAuthToken{
		AccessToken:  token.AccessToken,  // Токен доступа
		TokenType:    token.Type(),       // Тип токена
		RefreshToken: token.RefreshToken, // Токен обновления
		Expiry:       token.Expiry,       // Время истечения токена
		LinkID:       "",                 // Пустой ID, так как здесь неизвестно с каким пользователем будет связан токен
		Provider:     o.Name(),           // Имя провайдера
	}, nil
}

// User получает информацию о пользователе Google, используя токен OAuth.
func (o *Google) User(token domain.OAuthToken) (domain.OAuthUser, error) {
	const getUser = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="
	// Выполняет GET-запрос для получения информации о пользователе
	response, err := http.Get(getUser + token.AccessToken)
	if err != nil {
		return domain.OAuthUser{}, err // Возвращает ошибку, если запрос не удался
	}

	defer func() { _ = response.Body.Close() }() // Закрывает тело ответа после завершения работы с ним
	data, err := io.ReadAll(response.Body)       // Читает данные из тела ответа
	if err != nil {
		return domain.OAuthUser{}, err // Возвращает ошибку, если чтение не удалось
	}

	// Сложить данные в структуру ответа
	var googleUser struct {
		ID      string `json:"id"`      // Идентификатор пользователя
		Email   string `json:"email"`   // Электронная почта пользователя
		Name    string `json:"name"`    // Имя пользователя
		Picture string `json:"picture"` // URL изображения профиля пользователя
	}
	if err = json.Unmarshal(data, &googleUser); err != nil {
		return domain.OAuthUser{}, err // Возвращает ошибку, если разбор JSON не удался
	}

	return domain.OAuthUser{
		ID:       googleUser.ID,      // Идентификатор пользователя
		Email:    googleUser.Email,   // Электронная почта пользователя
		Name:     googleUser.Name,    // Имя пользователя
		Picture:  googleUser.Picture, // URL изображения профиля
		Provider: o.Name(),           // Имя провайдера
	}, nil
}

// AuthorizationURL генерирует URL для авторизации с использованием кода состояния.
func (o *Google) AuthorizationURL(state string) string {
	return o.config().AuthCodeURL(state) // Генерирует URL для авторизации
}
