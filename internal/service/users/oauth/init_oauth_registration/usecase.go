package initOAuthRegistration

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/service/users/oauth"
)

var (
	ErrInvalidProvider = errors.New("некорректное значение Provider")
)

// In представляет собой параметры инициализации регистрации OAuth.
type In struct {
	Provider string // Имя провайдера OAuth
}

// Validate валидирует значение параметра провайдера.
func (in In) Validate() error {
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	return nil // Возвращает nil, если параметры валидны
}

// Out представляет собой результат инициализации регистрации OAuth.
type Out struct {
	RedirectURL string // URL для перенаправления на страницу авторизации провайдера
}

type InitOAuthRegistrationUsecase struct {
	Providers oauth.OAuthProviders
}

// InitOAuthRegistration инициализирует процесс регистрации пользователя через OAuth.
func (u *InitOAuthRegistrationUsecase) InitOAuthRegistration(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Определить провайдера OAuth
	provider, err := u.Providers.Provider(in.Provider)
	if err != nil {
		return Out{}, err
	}

	// Генерирует URL для перенаправления на страницу авторизации провайдера
	return Out{
		RedirectURL: provider.AuthorizationURL(uuid.NewString()),
	}, nil
}
