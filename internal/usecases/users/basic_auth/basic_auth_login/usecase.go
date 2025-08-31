package basicAuthLogin

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

var (
	ErrInvalidProvider             = errors.New("некорректное значение Provider")
	ErrInvalidLogin                = errors.New("некорректное значение BasicAuthLogin")
	ErrInvalidPassword             = errors.New("некорректное значение Password")
	ErrLoginOrPasswordDoesNotMatch = errors.New("не совпадает BasicAuthLogin или Password")
)

type In struct {
	Login    string
	Password string
}

// Validate валидирует значение отдельно каждого параметры
func (in In) Validate() error {
	if err := userr.ValidateBasicAuthLogin(in.Login); err != nil {
		return errors.Join(ErrInvalidLogin, err)
	}
	if err := userr.ValidateBasicAuthPassword(in.Password); err != nil {
		return errors.Join(ErrInvalidPassword, err)
	}

	return nil
}

type Out struct {
	Session sessionn.Session
	User    userr.User
}

type BasicAuthLoginUsecase struct {
	Repo         userr.Repository
	SessionsRepo sessionn.Repository
}

func (u *BasicAuthLoginUsecase) BasicAuthLogin(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Получить метод входа
	matchUsers, err := u.Repo.List(userr.Filter{
		BasicAuthLogin:    in.Login,
		BasicAuthPassword: in.Password,
	})
	if err != nil {
		return Out{}, err
	}
	if len(matchUsers) != 1 {
		return Out{}, ErrLoginOrPasswordDoesNotMatch
	}
	user := matchUsers[0]

	// Создать сессию для пользователя
	sessionName := "todo: [название модели телефона / название браузера]"
	session, err := sessionn.NewSession(user.ID, sessionName, sessionn.StatusVerified)
	if err != nil {
		return Out{}, err
	}
	if err = u.SessionsRepo.Upsert(session); err != nil {
		return Out{}, err
	}

	return Out{
		Session: session,
		User:    user,
	}, nil
}
