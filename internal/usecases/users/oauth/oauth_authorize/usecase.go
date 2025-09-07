package oauthAuthorize

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

var (
	ErrInvalidProvider = errors.New("некорректное значение Provider")
)

// In представляет собой параметры инициализации регистрации Oauth.
type In struct {
	Provider string // Имя провайдера Oauth
}

// Validate валидирует значение параметра провайдера.
func (in In) Validate() error {
	if in.Provider == "" {
		return ErrInvalidProvider
	}
	return nil // Возвращает nil, если параметры валидны
}

// Out представляет собой результат инициализации регистрации Oauth.
type Out struct {
	RedirectURL string // URL для перенаправления на страницу авторизации провайдера
	State       string // Сгенерированная случайная строка для использования в предотвращении CSRF-атаки
}

type OauthAuthorizeUsecase struct {
	Providers oauth.Providers
}

// OauthAuthorize инициализирует процесс регистрации пользователя через Oauth.
func (u *OauthAuthorizeUsecase) OauthAuthorize(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Определить провайдера Oauth
	provider, err := u.Providers.Provider(in.Provider)
	if err != nil {
		return Out{}, err
	}

	state := uuid.NewString()
	return Out{
		// Генерирует URL для перенаправления на страницу авторизации провайдера
		RedirectURL: provider.AuthorizationURL(state),
		State:       state,
	}, nil
}
