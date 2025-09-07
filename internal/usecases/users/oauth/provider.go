package oauth

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/userr"
)

var (
	ErrUnknownOauthProvider = errors.New("неизвестный Oauth провайдер")
)

// Provider определяет интерфейс для работы с провайдерами Oauth.
type Provider interface {
	// Exchange обменивает код авторизации на токен Oauth
	Exchange(code string) (userr.OpenAuthToken, error)

	// User возвращает информацию о пользователе провайдера, используя токен Oauth
	User(token userr.OpenAuthToken) (userr.OpenAuthUser, error)

	// AuthorizationURL возвращает URL для авторизации.
	//
	// Параметр state используется для предотвращения CSRF-атаки.
	// Должен быть уникальной случайной строкой
	AuthorizationURL(state string) string

	// Name возвращает имя провайдера Oauth
	Name() string
}

// Providers представляет собой карту провайдеров Oauth, где ключом является имя провайдера.
type Providers map[string]Provider

func (o Providers) Add(p Provider) {
	o[p.Name()] = p
}

// Provider возвращает провайдера Oauth по его имени
func (o Providers) Provider(provider string) (Provider, error) {
	p, ok := o[provider]
	// Проверить, существует ли провайдер в карте
	if !ok || p == nil {
		return nil, ErrUnknownOauthProvider
	}

	return p, nil
}
