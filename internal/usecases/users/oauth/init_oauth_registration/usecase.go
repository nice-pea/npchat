package initOauthRegistration

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
}

type InitOauthRegistrationUsecase struct {
	Providers oauth.OauthProviders
}

// InitOauthRegistration инициализирует процесс регистрации пользователя через Oauth.
func (u *InitOauthRegistrationUsecase) InitOauthRegistration(in In) (Out, error) {
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
		RedirectURL: provider.AuthorizationURL(uuid.NewString()),
	}, nil
}
