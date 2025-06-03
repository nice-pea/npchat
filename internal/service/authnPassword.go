package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

type AuthnPassword struct {
	Repo         userr.Repository
	SessionsRepo domain.SessionsRepository
}

type AuthnPasswordLoginInput struct {
	Login    string
	Password string
}
type AuthnPasswordLoginOutput struct {
	Session domain.Session
	User    userr.User
}

// Validate валидирует значение отдельно каждого параметры
func (in AuthnPasswordLoginInput) Validate() error {
	if in.Login == "" {
		return errors.New("login is required")
	}
	if in.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (l *AuthnPassword) Login(in AuthnPasswordLoginInput) (AuthnPasswordLoginOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return AuthnPasswordLoginOutput{}, err
	}

	// Получить метод входа
	matchUsers, err := l.Repo.List(userr.Filter{
		BasicAuthLogin:    in.Login,
		BasicAuthPassword: in.Password,
	})
	if err != nil {
		return AuthnPasswordLoginOutput{}, err
	}
	if len(matchUsers) != 1 {
		return AuthnPasswordLoginOutput{}, ErrLoginOrPasswordDoesNotMatch
	}
	user := matchUsers[0]

	// Создать сессию для пользователя
	session := domain.Session{
		ID:     uuid.NewString(),
		UserID: user.ID,
		Token:  uuid.NewString(),
		Status: domain.SessionStatusVerified,
	}
	if err = l.SessionsRepo.Save(session); err != nil {
		return AuthnPasswordLoginOutput{}, err
	}

	return AuthnPasswordLoginOutput{
		Session: session,
		User:    user,
	}, nil
}

type AuthnPasswordRegistrationInput struct {
	Login    string
	Password string
	Name     string
	Nick     string
}

func (in AuthnPasswordRegistrationInput) Validate() error {
	if in.Login == "" {
		return errors.New("login is required")
	}
	if in.Password == "" {
		return errors.New("password is required")
	}
	if in.Name == "" {
		return errors.New("name is required")
	}

	return nil
}

type AuthnPasswordRegistrationOutput struct {
	Session domain.Session
	User    userr.User
}

func (l *AuthnPassword) Registration(in AuthnPasswordRegistrationInput) (AuthnPasswordRegistrationOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}

	// Создать метод аутентификации по логину и паролю
	basicAuth, err := userr.NewBasicAuth(in.Login, in.Password)
	if err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}

	// Создать пользователя
	user, err := userr.NewUser(in.Name, in.Nick)
	if err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}
	// Добавить метод аутентификации по логину и паролю
	if err := user.AddBasicAuth(basicAuth); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}

	// Проверка на существование пользователя с таким логином
	if conflictUsers, err := l.Repo.List(userr.Filter{
		BasicAuthLogin: in.Login,
	}); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	} else if len(conflictUsers) > 0 {
		return AuthnPasswordRegistrationOutput{}, ErrLoginIsAlreadyInUse
	}

	// Сохранить пользователя в репозиторий
	if err := l.Repo.Upsert(user); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}

	// Создать сессию для пользователя
	session := domain.Session2{
		UserID: user.ID,
		Name:   "todo: [название модели телефона / название браузера]",
		Status: domain.SessionStatusVerified, // Подтвержденная сессия
		AccessToken: domain.SessionToken{
			Token:  uuid.NewString(),
			Expiry: time.Now().Add(time.Minute * 10),
		},
		RefreshToken: domain.SessionToken{
			Token:  uuid.NewString(),
			Expiry: time.Now().Add(time.Hour * 24 * 60),
		},
	}
	if err := l.SessionsRepo.Save(session); err != nil {
		return AuthnPasswordRegistrationOutput{}, err
	}

	return AuthnPasswordRegistrationOutput{
		Session: session,
		User:    user,
	}, nil
}
