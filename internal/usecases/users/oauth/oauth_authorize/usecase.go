package oauthAuthorize

import (
	"errors"
	"net/url"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

var (
	ErrInvalidProvider         = errors.New("некорректное значение Provider")
	ErrInvalidCompleteCallback = errors.New("некорректное значение CompleteCallback")
)

// In представляет собой параметры инициализации регистрации Oauth.
type In struct {
	Provider         string // Имя провайдера Oauth
	CompleteCallback string // URL для перенаправления после авторизации
}

// Validate валидирует значение параметра провайдера.
func (in In) Validate() error {
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	u, err := url.Parse(in.CompleteCallback)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return ErrInvalidCompleteCallback
	}

	return nil // Возвращает nil, если параметры валидны
}

// Out представляет собой результат инициализации регистрации Oauth.
type Out struct {
	RedirectURL string // URL для перенаправления на страницу авторизации провайдера
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

	// Генерирует URL для перенаправления на страницу авторизации провайдера
	return Out{
		RedirectURL: provider.AuthorizationURL(uuid.NewString(), in.CompleteCallback),
	}, nil
}
