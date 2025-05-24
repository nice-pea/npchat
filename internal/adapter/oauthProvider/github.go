package oauthProvider

import (
	"context"
	"encoding/json"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// GitHub представляет собой структуру для работы с OAuth2 аутентификацией через GitHub.
type GitHub struct {
	ClientID     string // Идентификатор клиента для OAuth2
	ClientSecret string // Секрет клиента для OAuth2
	RedirectURL  string // URL для перенаправления после аутентификации
}

// Name возвращает имя провайдера OAuth.
func (o *GitHub) Name() string {
	return "github"
}

// config создает и возвращает конфигурацию OAuth2 для GitHub.
func (o *GitHub) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.ClientID,
		ClientSecret: o.ClientSecret,
		Endpoint:     github.Endpoint, // Использует конечную точку GitHub для OAuth2
		RedirectURL:  o.RedirectURL,
		Scopes:       []string{"user:email"}, // Запрашиваем доступ к электронной почте пользователя
	}
}

// Exchange обменивает код авторизации на токен OAuth.
func (o *GitHub) Exchange(code string) (domain.OAuthToken, error) {
	// Обменять кода авторизации на токен OAuth
	token, err := o.config().Exchange(context.Background(), code)
	if err != nil {
		return domain.OAuthToken{}, err // Возвращает ошибку, если обмен не удался
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

// User получает информацию о пользователе GitHub, используя токен OAuth.
func (o *GitHub) User(token domain.OAuthToken) (domain.OAuthUser, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken}, // Создает HTTP-клиент с токеном доступа
	))

	// Сделать запрос на GitHub для получения информации о пользователе
	const getUser = "https://api.github.com/user"
	response, err := client.Get(getUser)
	if err != nil {
		return domain.OAuthUser{}, err
	}
	defer func() { _ = response.Body.Close() }()

	// Прочитать данные из тела ответа
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return domain.OAuthUser{}, err
	}

	// Сложить данные в структуру ответа
	var githubUser struct {
		ID        int64  `json:"id"`         // Идентификатор пользователя
		Login     string `json:"login"`      // Логин пользователя
		Name      string `json:"name"`       // Имя пользователя
		Email     string `json:"email"`      // Электронная почта пользователя
		AvatarURL string `json:"avatar_url"` // URL изображения профиля пользователя
	}
	if err := json.Unmarshal(data, &githubUser); err != nil {
		return domain.OAuthUser{}, err
	}

	return domain.OAuthUser{
		ID:       string(rune(githubUser.ID)),
		Email:    githubUser.Email,
		Name:     githubUser.Name,
		Picture:  githubUser.AvatarURL,
		Provider: o.Name(),
	}, nil
}

// AuthorizationURL генерирует URL для авторизации с использованием кода состояния.
func (o *GitHub) AuthorizationURL(state string) string {
	return o.config().AuthCodeURL(state)
}
