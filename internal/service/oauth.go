package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type OAuth struct {
	Providers adapter.OAuthProviders
	OAuthRepo domain.OAuthRepository
	UsersRepo domain.UsersRepository
}

type OAuthCompeteRegistrationInput struct {
	UserCode  string
	InitState string
	Provider  string
}

// Validate валидирует значение отдельно каждого параметры
func (in OAuthCompeteRegistrationInput) Validate() error {
	if in.UserCode == "" {
		return ErrInvalidUserCode
	}
	if in.InitState == "" {
		return ErrInvalidInitState
	}
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	return nil
}

// CompeteRegistration
// Подсмотрено в: https://github.com/oguzhantasimaz/Go-Clean-Architecture-Template/blob/main/api/controller/google.go
func (o *OAuth) CompeteRegistration(in OAuthCompeteRegistrationInput) (domain.User, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return domain.User{}, err
	}

	// Определить провайдера OAuth
	provider, err := o.provider(in.Provider)
	if err != nil {
		return domain.User{}, err
	}

	// Найти ранее созданную пустую связь InitState
	links, err := o.OAuthRepo.ListLinks(domain.OAuthListLinksFilter{
		ID:       in.InitState,
		Provider: in.Provider,
	})
	if err != nil {
		return domain.User{}, err
	}
	if len(links) != 1 {
		return domain.User{}, ErrWrongInitState
	}
	if links[0].UserID != "" || links[0].ExternalID != "" {
		return domain.User{}, ErrOAuthRegistrationAlreadyCompleted
	}
	oauthLink := links[0]

	// Получить пользователя провайдера
	token, err := provider.Exchange(in.UserCode)
	if err != nil {
		return domain.User{}, errors.Join(ErrWrongUserCode, err)
	}
	pUser, err := provider.User(token)
	if err != nil {
		return domain.User{}, err
	}

	// Проверить, не связан ли пользователь провайдера с каким-нибудь нашим пользователем
	links, err = o.OAuthRepo.ListLinks(domain.OAuthListLinksFilter{
		ExternalID: pUser.ID,
		Provider:   in.Provider,
	})
	if err != nil {
		return domain.User{}, err
	}
	if len(links) != 0 {
		return domain.User{}, ErrProvidersUserIsAlreadyLinked
	}

	// Создать пользователя
	user := domain.User{
		ID:   uuid.NewString(),
		Name: pUser.Name,
		Nick: "",
	}
	if err = o.UsersRepo.Save(user); err != nil {
		return domain.User{}, err
	}

	// Сохранить связь нашего пользователя с пользователем провайдера
	err = o.OAuthRepo.SaveLink(domain.OAuthLink{
		ID:         oauthLink.ID,
		UserID:     user.ID,
		ExternalID: pUser.ID,
		Provider:   in.Provider,
	})
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

type OAuthRegistrationInitOut struct {
	RedirectURL string
}

type OAuthInitRegistrationInput struct {
	Provider string
}

func (in OAuthInitRegistrationInput) Validate() error {
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	return nil
}

func (o *OAuth) InitRegistration(in OAuthInitRegistrationInput) (OAuthRegistrationInitOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return OAuthRegistrationInitOut{}, err
	}

	// Определить провайдера OAuth
	provider, err := o.provider(in.Provider)
	if err != nil {
		return OAuthRegistrationInitOut{}, err
	}

	// Сохранить связь
	link := domain.OAuthLink{
		ID:         uuid.NewString(),
		UserID:     "",
		ExternalID: "",
		Provider:   in.Provider,
	}
	if err := o.OAuthRepo.SaveLink(link); err != nil {
		return OAuthRegistrationInitOut{}, err
	}

	return OAuthRegistrationInitOut{
		RedirectURL: provider.AuthCodeURL(link.ID),
	}, nil
}

func (o *OAuth) provider(provider string) (adapter.OAuthProvider, error) {
	p, ok := o.Providers[provider]
	if !ok || p == nil {
		return nil, ErrUnknownOAuthProvider
	}

	return p, nil
}
