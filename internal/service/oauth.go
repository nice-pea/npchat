package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type OAuth struct {
	Providers    OAuthProviders
	OAuthRepo    domain.OAuthRepository
	UsersRepo    domain.UsersRepository
	SessionsRepo domain.SessionsRepository
}

type OAuthProviders = map[string]OAuthProvider

type OAuthProvider interface {
	Exchange(code string) (domain.OAuthToken, error)
	User(token domain.OAuthToken) (domain.OAuthUser, error)
	AuthCodeURL(state string) string
}

type OAuthCompeteRegistrationInput struct {
	UserCode string
	Provider string
}

// Validate валидирует значение отдельно каждого параметры
func (in OAuthCompeteRegistrationInput) Validate() error {
	if in.UserCode == "" {
		return ErrInvalidUserCode
	}
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	return nil
}

type OAuthCompeteRegistrationOut struct {
	Session domain.Session
	User    domain.User
}

// CompeteRegistration
// Подсмотрено в: https://github.com/oguzhantasimaz/Go-Clean-Architecture-Template/blob/main/api/controller/google.go
func (o *OAuth) CompeteRegistration(in OAuthCompeteRegistrationInput) (OAuthCompeteRegistrationOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Определить провайдера OAuth
	provider, err := o.provider(in.Provider)
	if err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Получить пользователя провайдера
	token, err := provider.Exchange(in.UserCode)
	if err != nil {
		return OAuthCompeteRegistrationOut{}, errors.Join(ErrWrongUserCode, err)
	}
	pUser, err := provider.User(token)
	if err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Проверить, не связан ли пользователь провайдера с каким-нибудь нашим пользователем
	links, err := o.OAuthRepo.ListLinks(domain.OAuthListLinksFilter{
		ExternalID: pUser.ID,
		Provider:   in.Provider,
	})
	if err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}
	if len(links) != 0 {
		return OAuthCompeteRegistrationOut{}, ErrProvidersUserIsAlreadyLinked
	}

	// Создать пользователя
	user := domain.User{
		ID:   uuid.NewString(),
		Name: pUser.Name,
		Nick: "",
	}
	if err = o.UsersRepo.Save(user); err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Сохранить связь нашего пользователя с пользователем провайдера
	err = o.OAuthRepo.SaveLink(domain.OAuthLink{
		UserID:     user.ID,
		ExternalID: pUser.ID,
		Provider:   in.Provider,
	})
	if err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Создать сессию для пользователя
	session := domain.Session{
		ID:     uuid.NewString(),
		UserID: user.ID,
		Token:  uuid.NewString(),
		Status: domain.SessionStatusVerified,
	}
	if err = o.SessionsRepo.Save(session); err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	return OAuthCompeteRegistrationOut{
		Session: session,
		User:    user,
	}, nil
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

	return OAuthRegistrationInitOut{
		RedirectURL: provider.AuthCodeURL(uuid.NewString()),
	}, nil
}

func (o *OAuth) provider(provider string) (OAuthProvider, error) {
	p, ok := o.Providers[provider]
	if !ok || p == nil {
		return nil, ErrUnknownOAuthProvider
	}

	return p, nil
}

type OAuthInitLoginInput struct {
	Provider string
}

type OAuthInitLoginOut struct {
	RedirectURL string
}

func (o *OAuth) InitLogin(in OAuthInitLoginInput) (OAuthInitLoginOut, error) {
	return OAuthInitLoginOut{}, nil
}

type OAuthCompleteLoginInput struct {
	UserCode  string
	InitState string
	Provider  string
}

type OAuthCompleteLoginOut struct {
	Session domain.Session
	User    domain.User
}

func (o *OAuth) CompleteLogin(in OAuthCompleteLoginInput) (OAuthCompleteLoginOut, error) {
	return OAuthCompleteLoginOut{}, nil
}
