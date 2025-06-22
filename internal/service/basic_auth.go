package service

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

type BasicAuthLoginIn struct {
	Login    string
	Password string
}
type BasicAuthLoginOut struct {
	Session sessionn.Session
	User    userr.User
}

// Validate валидирует значение отдельно каждого параметры
func (in BasicAuthLoginIn) Validate() error {
	if err := userr.ValidateBasicAuthLogin(in.Login); err != nil {
		return err
	}
	if err := userr.ValidateBasicAuthPassword(in.Password); err != nil {
		return err
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
	session, err := sessionn.NewSession(user.ID, "", sessionn.StatusVerified)
	if err != nil {
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
	Session sessionn.Session
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
	session, err := sessionn.NewSession(user.ID, "", sessionn.StatusVerified)
	if err != nil {
		return BasicAuthRegistrationOut{}, err
	}

	return BasicAuthRegistrationOut{
		Session: session,
		User:    user,
	}, nil
}
