package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type AuthnPassword struct {
	AuthnPasswordRepo domain.AuthnPasswordRepository
	SessionsRepo      domain.SessionsRepository
}
type AuthnPasswordLoginInput struct {
	Login    string
	Password string
}

// Validate валидирует значение отдельно каждого параметры
func (in AuthnPasswordLoginInput) Validate() error {
	lc := domain.AuthnPassword{
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

func (l *AuthnPassword) Login(in AuthnPasswordLoginInput) (domain.Session, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return domain.Session{}, err
	}

	// Получить пользователя по логину и паролю
	aps, err := l.AuthnPasswordRepo.List(domain.AuthnPasswordFilter{
		Login:    in.Login,
		Password: in.Password,
	})
	if err != nil {
		return domain.Session{}, err
	}
	if len(aps) != 1 {
		return domain.Session{}, ErrLoginOrPasswordDoesNotMatch
	}

	// Создать сессию для пользователя
	session := domain.Session{
		ID:     uuid.NewString(),
		UserID: aps[0].UserID,
		Token:  uuid.NewString(),
		Status: domain.SessionStatusVerified,
	}
	if err = l.SessionsRepo.Save(session); err != nil {
		return domain.Session{}, err
	}

	return session, nil
}
