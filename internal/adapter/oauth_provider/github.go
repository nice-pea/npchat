package oauth_provider

import (
	"context"
	"encoding/json"
	"io"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

// GitHub представляет собой структуру для работы с OAuth2 аутентификацией через GitHub.
type GitHub struct {
	clientID     string // Идентификатор клиента для OAuth2
	clientSecret string // Секрет клиента для OAuth2
	redirectURL  string // URL для перенаправления после аутентификации
}

type GitHubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

func NewGitHub(cfg GitHubConfig) *GitHub {
	return &GitHub{
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		redirectURL:  cfg.RedirectURL,
	}
}

// Name возвращает имя провайдера OAuth.
func (o *GitHub) Name() string {
	return "github"
}

// config создает и возвращает конфигурацию OAuth2 для GitHub.
func (o *GitHub) config() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     o.clientID,
		ClientSecret: o.clientSecret,
		Endpoint:     github.Endpoint, // Использует конечную точку GitHub для OAuth2
		RedirectURL:  o.redirectURL,
		Scopes:       []string{"user:email"}, // Запрашиваем доступ к электронной почте пользователя
	}
}

// Exchange обменивает код авторизации на токен OAuth.
func (o *GitHub) Exchange(code string) (userr.OpenAuthToken, error) {
	// Обменять кода авторизации на токен OAuth
	token, err := o.config().Exchange(context.Background(), code)
	if err != nil {
		return userr.OpenAuthToken{}, err // Возвращает ошибку, если обмен не удался
	}

	return userr.NewOpenAuthToken(
		token.AccessToken,  // Токен доступа
		token.Type(),       // Тип токена
		token.RefreshToken, // Токен обновления
		token.Expiry,       // Время истечения токена
	)
}

// User получает информацию о пользователе GitHub, используя токен OAuth.
func (o *GitHub) User(token userr.OpenAuthToken) (userr.OpenAuthUser, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken}, // Создает HTTP-клиент с токеном доступа
	))

	// Сделать запрос на GitHub для получения информации о пользователе
	const getUser = "https://api.github.com/user"
	response, err := client.Get(getUser)
	if err != nil {
		return userr.OpenAuthUser{}, err
	}
	defer func() { _ = response.Body.Close() }()

	// Прочитать данные из тела ответа
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return userr.OpenAuthUser{}, err
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
		return userr.OpenAuthUser{}, err
	}

	return userr.NewOpenAuthUser(
		string(rune(githubUser.ID)), // Идентификатор пользователя
		o.Name(),                    // Имя провайдера
		githubUser.Email,            // Электронная почта пользователя
		githubUser.Name,             // Имя пользователя
		githubUser.AvatarURL,        // URL изображения профиля
		token,
	)
}

// AuthorizationURL генерирует URL для авторизации с использованием кода состояния.
func (o *GitHub) AuthorizationURL(state string) string {
	return o.config().AuthCodeURL(state)
}
