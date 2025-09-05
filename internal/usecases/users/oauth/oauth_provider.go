package oauth

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/userr"
)

var (
	ErrUnknownOauthProvider = errors.New("неизвестный Oauth провайдер")
)

// OauthProvider определяет интерфейс для работы с провайдерами Oauth.
type OauthProvider interface {
	// Exchange обменивает код авторизации на токен Oauth
	Exchange(code string) (userr.OpenAuthToken, error)

	// User возвращает информацию о пользователе провайдера, используя токен Oauth
	User(token userr.OpenAuthToken) (userr.OpenAuthUser, error)

	// AuthorizationURL возвращает URL для авторизации.
	// Параметр state используется для предотвращения CSRF-атаки, Должен быть уникальной случайной строкой
	AuthorizationURL(state string) string

	// Name возвращает имя провайдера Oauth
	Name() string
}

// OauthProviders представляет собой карту провайдеров Oauth, где ключом является имя провайдера.
type OauthProviders map[string]OauthProvider

func (o OauthProviders) Add(p OauthProvider) {
	o[p.Name()] = p
}

// provider возвращает провайдера Oauth по его имени
func (o OauthProviders) Provider(provider string) (OauthProvider, error) {
	p, ok := o[provider]
	// Проверить, существует ли провайдер в карте
	if !ok || p == nil {
		return nil, ErrUnknownOauthProvider
	}

	return p, nil
}
