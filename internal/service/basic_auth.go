package service

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

type BasicAuthLoginIn struct {
	Login    string
	Password string
}
type BasicAuthLoginOut struct {
	Session domain.Session
	User    userr.User
}

// Validate валидирует значение отдельно каждого параметры
func (in BasicAuthLoginIn) Validate() error {
	if in.Login == "" {
		return errors.New("login is required")
	}
	if in.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

func (u *Users) BasicAuthLogin(in BasicAuthLoginIn) (BasicAuthLoginOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return BasicAuthLoginOut{}, err
	}

	// Получить метод входа
	matchUsers, err := u.Repo.List(userr.Filter{
		BasicAuthLogin:    in.Login,
		BasicAuthPassword: in.Password,
	})
	if err != nil {
		return BasicAuthLoginOut{}, err
	}
	if len(matchUsers) != 1 {
		return BasicAuthLoginOut{}, ErrLoginOrPasswordDoesNotMatch
	}
	user := matchUsers[0]

	// Создать сессию для пользователя
	session := domain.Session{
		ID:     uuid.NewString(),
		UserID: user.ID,
		Token:  uuid.NewString(),
		Status: domain.SessionStatusVerified,
	}
	if err = u.SessionsRepo.Save(session); err != nil {
		return BasicAuthLoginOut{}, err
	}

	return BasicAuthLoginOut{
		Session: session,
		User:    user,
	}, nil
}

type BasicAuthRegistrationIn struct {
	Login    string
	Password string
	Name     string
	Nick     string
}

func (in BasicAuthRegistrationIn) Validate() error {
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

type BasicAuthRegistrationOut struct {
	Session domain.Session
	User    userr.User
}

func (u *Users) BasicAuthRegistration(in BasicAuthRegistrationIn) (BasicAuthRegistrationOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return BasicAuthRegistrationOut{}, err
	}

	// Создать метод аутентификации по логину и паролю
	basicAuth, err := userr.NewBasicAuth(in.Login, in.Password)
	if err != nil {
		return BasicAuthRegistrationOut{}, err
	}

	// Создать пользователя
	user, err := userr.NewUser(in.Name, in.Nick)
	if err != nil {
		return BasicAuthRegistrationOut{}, err
	}
	// Добавить метод аутентификации по логину и паролю
	if err := user.AddBasicAuth(basicAuth); err != nil {
		return BasicAuthRegistrationOut{}, err
	}

	// Проверка на существование пользователя с таким логином
	if conflictUsers, err := u.Repo.List(userr.Filter{
		BasicAuthLogin: in.Login,
	}); err != nil {
		return BasicAuthRegistrationOut{}, err
	} else if len(conflictUsers) > 0 {
		return BasicAuthRegistrationOut{}, ErrLoginIsAlreadyInUse
	}

	// Сохранить пользователя в репозиторий
	if err := u.Repo.Upsert(user); err != nil {
		return BasicAuthRegistrationOut{}, err
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
	if err := u.SessionsRepo.Save(session); err != nil {
		return BasicAuthRegistrationOut{}, err
	}

	return BasicAuthRegistrationOut{
		Session: session,
		User:    user,
	}, nil
}
