package oauth

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/userr"
)

var (
	ErrUnknownOAuthProvider = errors.New("неизвестный OAuth провайдер")
)

// OAuthProvider определяет интерфейс для работы с провайдерами OAuth.
type OAuthProvider interface {
	// Exchange обменивает код авторизации на токен OAuth
	Exchange(code string) (userr.OpenAuthToken, error)

	// User возвращает информацию о пользователе провайдера, используя токен OAuth
	User(token userr.OpenAuthToken) (userr.OpenAuthUser, error)

	// AuthorizationURL возвращает URL для авторизации.
	// Параметр state используется для предотвращения CSRF-атаки, Должен быть уникальной случайной строкой
	AuthorizationURL(state string) string

	// Name возвращает имя провайдера OAuth
	Name() string
}

// OAuthProviders представляет собой карту провайдеров OAuth, где ключом является имя провайдера.
type OAuthProviders map[string]OAuthProvider

func (o OAuthProviders) Add(p OAuthProvider) {
	o[p.Name()] = p
}

// provider возвращает провайдера OAuth по его имени
func (o OAuthProviders) Provider(provider string) (OAuthProvider, error) {
	p, ok := o[provider]
	// Проверить, существует ли провайдер в карте
	if !ok || p == nil {
		return nil, ErrUnknownOAuthProvider
	}

	return p, nil
}
