package oauth_provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"

	"github.com/nice-pea/npchat/internal/domain/userr"
)

// GitHub представляет собой структуру для работы с OAuth2 аутентификацией через GitHub.
type GitHub struct {
	config *oauth2.Config // Конфигурация OAuth2 для GitHub
}

type GitHubConfig struct {
	ClientID     string // Идентификатор клиента для OAuth2
	ClientSecret string // Секрет клиента для OAuth2
	RedirectURL  string // URL для перенаправления после аутентификации
}

func NewGitHub(cfg GitHubConfig) *GitHub {
	return &GitHub{
		config: &oauth2.Config{
			ClientID:     cfg.ClientID,
			ClientSecret: cfg.ClientSecret,
			Endpoint:     github.Endpoint, // Использует конечную точку GitHub для OAuth2
			RedirectURL:  cfg.RedirectURL,
			Scopes:       []string{"user:email"}, // Запрашиваем доступ к электронной почте пользователя
		},
	}
}

// Name возвращает имя провайдера OAuth.
func (o *GitHub) Name() string {
	return "github"
}

// Exchange обменивает код авторизации на токен OAuth.
func (o *GitHub) Exchange(code string) (userr.OpenAuthToken, error) {
	// Обменять кода авторизации на токен OAuth
	token, err := o.config.Exchange(context.Background(), code)
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

const githubGetUserUrl = "https://api.github.com/user"

// User получает информацию о пользователе GitHub, используя токен OAuth.
func (o *GitHub) User(token userr.OpenAuthToken) (userr.OpenAuthUser, error) {
	client := oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token.AccessToken}, // Создает HTTP-клиент с токеном доступа
	))

	// Сделать запрос на GitHub для получения информации о пользователе
	response, err := client.Get(githubGetUserUrl)
	if err != nil {
		return userr.OpenAuthUser{}, err
	}
	defer func() { _ = response.Body.Close() }()

	// Проверить код ответа
	if response.StatusCode != http.StatusOK {
		return userr.OpenAuthUser{}, fmt.Errorf("unexpected status code: %d", response.StatusCode)
	}

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
		strconv.FormatInt(githubUser.ID, 10), // Идентификатор пользователя
		o.Name(),                             // Имя провайдера
		githubUser.Email,                     // Электронная почта пользователя
		githubUser.Name,                      // Имя пользователя
		githubUser.AvatarURL,                 // URL изображения профиля
		token,
	)
}

// AuthorizationURL генерирует URL для авторизации с использованием кода состояния.
func (o *GitHub) AuthorizationURL(state string) string {
	return o.config.AuthCodeURL(state)
}
