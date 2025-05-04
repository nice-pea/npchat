package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type LoginCredentials struct {
	LoginCredentialsRepo domain.LoginCredentialsRepository
	SessionsRepo         domain.SessionsRepository
}
type LoginByCredentialsInput struct {
	Login    string
	Password string
}

// Validate валидирует значение отдельно каждого параметры
func (in LoginByCredentialsInput) Validate() error {
	lc := domain.LoginCredentials{
		UserID:   "",
		Login:    in.Login,
		Password: in.Password,
	}
	if err := lc.ValidateLogin(); err != nil {
		return errors.Join(err, ErrInvalidLogin)
	}

	//if err := lc.ValidatePassword(); err != nil {
	//	return errors.Join(err, ErrInvalidPassword)
	//}

	return nil
}

var ErrLoginOrPasswordDoesNotMatch = errors.New("не совпадает Login или Password")

func (l *LoginCredentials) Login(in LoginByCredentialsInput) (domain.Session, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return domain.Session{}, err
	}

	// Получить пользователя по логину и паролю
	creds, err := l.LoginCredentialsRepo.List(domain.LoginCredentialsFilter{
		Login:    in.Login,
		Password: in.Password,
	})
	if err != nil {
		return domain.Session{}, err
	}
	if len(creds) != 1 {
		return domain.Session{}, ErrLoginOrPasswordDoesNotMatch
	}

	// Создать сессию для пользователя
	session := domain.Session{
		ID:     uuid.NewString(),
		UserID: creds[0].UserID,
		Token:  uuid.NewString(),
		Status: domain.SessionStatusVerified,
	}
	if err = l.SessionsRepo.Save(session); err != nil {
		return domain.Session{}, err
	}

	return session, nil
}
