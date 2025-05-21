package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/adapter"
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type OAuth struct {
	Google    adapter.OAuthGoogle
	OAuthRepo domain.OAuthRepository
	UsersRepo domain.UsersRepository
}

type OAuthRegistrationCallbackInput struct {
	Provider string
}

type GoogleRegistrationInput struct {
	UserCode  string
	InitState string
}

//type GoogleRegistrationOut struct {
//	User    domain.User
//	Session domain.Session
//}

// Validate валидирует значение отдельно каждого параметры
func (in GoogleRegistrationInput) Validate() error {
	if in.UserCode == "" {
		return ErrInvalidUserCode
	}
	if in.InitState == "" {
		return ErrInvalidInitState
	}

	return nil
}

// GoogleRegistration
// Подсмотрено в: https://github.com/oguzhantasimaz/Go-Clean-Architecture-Template/blob/main/api/controller/google.go
func (o *OAuth) GoogleRegistration(in GoogleRegistrationInput) (domain.User, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return domain.User{}, err
	}

	// Проверить InitState
	links, err := o.OAuthRepo.ListLinks(domain.OAuthListLinksFilter{ID: in.InitState})
	if err != nil {
		return domain.User{}, err
	}
	if len(links) != 1 {
		return domain.User{}, ErrWrongInitState
	}

	// Получить пользователя google
	token, err := o.Google.Exchange(in.UserCode)
	if err != nil {
		return domain.User{}, errors.Join(ErrWrongUserCode, err)
	}
	googleUser, err := o.Google.User(token)
	if err != nil {
		return domain.User{}, err
	}

	// Создать пользователя
	user := domain.User{
		ID:   uuid.NewString(),
		Name: googleUser.Name,
		Nick: "",
	}
	if err = o.UsersRepo.Save(user); err != nil {
		return domain.User{}, err
	}

	// Связать пользователя с google пользователем
	err = o.OAuthRepo.SaveLink(domain.OAuthLink{
		ID:         links[0].ID,
		UserID:     user.ID,
		ExternalID: googleUser.ID,
	})
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

type GoogleRegistrationInitOut struct {
	RedirectURL string
}

func (o *OAuth) GoogleRegistrationInit() (GoogleRegistrationInitOut, error) {
	// Сохранить связь
	link := domain.OAuthLink{
		ID:         uuid.NewString(),
		UserID:     "",
		ExternalID: "",
	}
	if err := o.OAuthRepo.SaveLink(link); err != nil {
		return GoogleRegistrationInitOut{}, err
	}

	return GoogleRegistrationInitOut{
		RedirectURL: o.Google.AuthCodeURL(link.ID),
	}, nil
}
