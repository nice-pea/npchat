package basicAuthRegistration

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

var (
	ErrLoginIsRequired     = errors.New("login это обязательный параметр")
	ErrPasswordIsRequired  = errors.New("password это обязательный параметр")
	ErrNameIsRequired      = errors.New("name это обязательный параметр")
	ErrLoginIsAlreadyInUse = errors.New("логин уже используется")
)

type In struct {
	Login    string
	Password string
	Name     string
	Nick     string
}

func (in In) Validate() error {
	if in.Login == "" {
		return ErrLoginIsRequired
	}
	if in.Password == "" {
		return ErrPasswordIsRequired
	}
	if in.Name == "" {
		return ErrNameIsRequired
	}

	return nil
}

type Out struct {
	Session sessionn.Session
	User    userr.User
}

type BasicAuthRegistrationUsecase struct {
	Repo         userr.Repository
	SessionsRepo sessionn.Repository
}

func (u *BasicAuthRegistrationUsecase) BasicAuthRegistration(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Создать метод аутентификации по логину и паролю
	basicAuth, err := userr.NewBasicAuth(in.Login, in.Password)
	if err != nil {
		return Out{}, err
	}

	// Создать пользователя
	user, err := userr.NewUser(in.Name, in.Nick)
	if err != nil {
		return Out{}, err
	}
	// Добавить метод аутентификации по логину и паролю
	if err := user.AddBasicAuth(basicAuth); err != nil {
		return Out{}, err
	}

	// Обернуть работу с репозиторием в транзакцию
	var session sessionn.Session
	if err = u.Repo.InTransaction(func(txRepo userr.Repository) error {
		// Проверка на существование пользователя с таким логином
		if conflictUsers, err := u.Repo.List(userr.Filter{
			BasicAuthLogin: in.Login,
		}); err != nil {
			return err
		} else if len(conflictUsers) > 0 {
			return ErrLoginIsAlreadyInUse
		}

		// Сохранить пользователя в репозиторий
		if err = u.Repo.Upsert(user); err != nil {
			return err
		}

		// Создать сессию для пользователя
		sessionName := "todo: [название модели телефона / название браузера]"
		if session, err = sessionn.NewSession(user.ID, sessionName, sessionn.StatusVerified); err != nil {
			return err
		}
		if err = u.SessionsRepo.Upsert(session); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return Out{}, err
	}

	return Out{
		Session: session,
		User:    user,
	}, nil
}
