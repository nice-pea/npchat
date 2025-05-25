package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type AuthnPassword struct {
	AuthnPasswordRepo domain.AuthnPasswordRepository
	SessionsRepo      domain.SessionsRepository
	UsersRepo         domain.UsersRepository
}

type AuthnPasswordLoginInput struct {
	Login    string
	Password string
}
type AuthnPasswordLoginOutput struct {
	Session domain.Session
	User    domain.User
}

// Validate валидирует значение отдельно каждого параметры
func (in AuthnPasswordLoginInput) Validate() error {
	lc := domain.AuthnPassword{
		UserID:   "",
		Login:    in.Login,
		Password: in.Password,
	}
	if err := lc.ValidateLogin(); err != nil {
		return errors.Join(ErrInvalidLogin, err)
	}

	if err := lc.ValidatePassword(); err != nil {
		return errors.Join(ErrInvalidPassword, err)
	}

	return nil
}

func (l *AuthnPassword) Login(in AuthnPasswordLoginInput) (AuthnPasswordLoginOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return AuthnPasswordLoginOutput{}, err
	}

	// Получить метод входа
	aps, err := l.AuthnPasswordRepo.List(domain.AuthnPasswordFilter{
		Login:    in.Login,
		Password: in.Password,
	})
	if err != nil {
		return AuthnPasswordLoginOutput{}, err
	}
	if len(aps) != 1 {
		return AuthnPasswordLoginOutput{}, ErrLoginOrPasswordDoesNotMatch
	}

	// Получить пользователя
	users, err := l.UsersRepo.List(domain.UsersFilter{ID: aps[0].UserID})
	if err != nil {
		return AuthnPasswordLoginOutput{}, err
	}
	if len(users) != 1 {
		return AuthnPasswordLoginOutput{}, ErrUserNotExists
	}

	// Создать сессию для пользователя
	session := domain.Session{
		ID:     uuid.NewString(),
		UserID: aps[0].UserID,
		Token:  uuid.NewString(),
		Status: domain.SessionStatusVerified,
	}
	if err = l.SessionsRepo.Save(session); err != nil {
		return AuthnPasswordLoginOutput{}, err
	}

	return AuthnPasswordLoginOutput{
		Session: session,
		User:    users[0],
	}, nil
}

type AuthnPasswordRegistrationInput struct {
	Login    string
	Password string
	Name     string
	Nick     string
}

func (in AuthnPasswordRegistrationInput) Validate() error {
	// Валидация полей для метода аутентификации
	lc := domain.AuthnPassword{
		Login:    in.Login,
		Password: in.Password,
	}
	if err := lc.ValidateLogin(); err != nil {
		return errors.Join(err, ErrInvalidLogin)
	}
	if err := lc.ValidatePassword(); err != nil {
		return errors.Join(err, ErrInvalidPassword)
	}
	// Валидация полей для создания пользователя
	u := domain.User{Name: in.Name, Nick: in.Nick}
	if err := u.ValidateName(); err != nil {
		return errors.Join(err, ErrInvalidName)
	}
	if err := u.ValidateNick(); err != nil {
		return errors.Join(err, ErrInvalidNick)
	}

	return nil
}

type AuthnPasswordRegistrationOutput struct {
	Session domain.Session
	User    domain.User
}

func (l *AuthnPassword) Registration(in AuthnPasswordRegistrationInput) (AuthnPasswordRegistrationOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}

	// Проверка на существование пользователя с таким логином
	apFilter := domain.AuthnPasswordFilter{
		Login: in.Login,
	}
	if aps, err := l.AuthnPasswordRepo.List(apFilter); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	} else if len(aps) > 0 {
		return AuthnPasswordRegistrationOutput{}, ErrLoginIsAlreadyInUse
	}

	// Создать пользователя
	user := domain.User{
		ID:   uuid.NewString(),
		Name: in.Name,
		Nick: in.Nick,
	}
	if err := l.UsersRepo.Save(user); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}

	// Создать метод входа
	ap := domain.AuthnPassword{
		UserID:   user.ID,
		Login:    in.Login,
		Password: in.Password,
	}
	if err := l.AuthnPasswordRepo.Save(ap); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}

	// Создать сессию для пользователя
	session := domain.Session{
		ID:     uuid.NewString(),
		UserID: user.ID,
		Token:  uuid.NewString(),
		Status: domain.SessionStatusVerified, // Подтвержденная сессия
	}
	if err := l.SessionsRepo.Save(session); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}

	return AuthnPasswordRegistrationOutput{
		Session: session,
		User:    user,
	}, nil
}
