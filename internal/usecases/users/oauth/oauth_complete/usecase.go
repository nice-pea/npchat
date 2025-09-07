package oauthComplete

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	"github.com/nice-pea/npchat/internal/usecases/users/oauth"
)

var (
	ErrInvalidProvider              = errors.New("некорректное значение Provider")
	ErrProvidersUserIsAlreadyLinked = errors.New("пользователь Oauth-провайдера уже связан с пользователем")
	ErrInvalidUserCode              = errors.New("некорректное значение UserCode")
	ErrWrongUserCode                = errors.New("неправильный UserCode")
)

// In представляет собой структуру для входных данных завершения регистрации Oauth.
type In struct {
	UserCode string // Код пользователя, полученный от провайдера
	Provider string // Имя провайдера Oauth
}

// Validate валидирует значение каждого параметра.
func (in In) Validate() error {
	if in.UserCode == "" {
		return ErrInvalidUserCode
	}
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	return nil // Возвращает nil, если все параметры валидны
}

// Out представляет собой результат завершения регистрации Oauth.
type Out struct {
	Session sessionn.Session // Сессия пользователя
	User    userr.User       // Пользователь
}

type OauthCompleteUsecase struct {
	Repo         userr.Repository
	Providers    oauth.Providers
	SessionsRepo sessionn.Repository // Репозиторий сессий пользователей
}

// OauthComplete завершает процесс регистрации пользователя через Oauth.
func (u *OauthCompleteUsecase) OauthComplete(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Определить провайдера Oauth
	provider, err := u.Providers.Provider(in.Provider)
	if err != nil {
		return Out{}, err
	}
	// Получить пользователя провайдера
	token, err := provider.Exchange(in.UserCode)
	if err != nil {
		return Out{}, errors.Join(ErrWrongUserCode, err)
	}
	openAuthUser, err := provider.User(token)
	if err != nil {
		return Out{}, err
	}

	// Определить, надо ли создавать пользователя
	user, userIsExists, err := u.userIfExists(openAuthUser)
	if err != nil {
		return Out{}, err
	}

	// Зарегистрировать либо авторизовать пользователя
	if userIsExists {
		return u.login(user)
	} else {
		return u.registration(openAuthUser)
	}
}

func (u *OauthCompleteUsecase) userIfExists(openAuthUser userr.OpenAuthUser) (userr.User, bool, error) {
	user, err := userr.Find(u.Repo, userr.Filter{
		OauthUserID:   openAuthUser.ID,
		OauthProvider: openAuthUser.Provider,
	})
	if errors.Is(err, userr.ErrUserNotExists) {
		return userr.User{}, false, nil
	} else if err != nil {
		return userr.User{}, false, err
	}

	return user, true, nil
}

func (u *OauthCompleteUsecase) registration(openAuthUser userr.OpenAuthUser) (Out, error) {
	// Создать пользователя
	user, err := userr.NewUser(openAuthUser.Name, "")
	if err != nil {
		return Out{}, err
	}
	// Добавить пользователя в список связей для аутентификации по Oauth
	if err = user.AddOpenAuthUser(openAuthUser); err != nil {
		return Out{}, err
	}

	// Обернуть работу с репозиторием в транзакцию
	var session sessionn.Session
	if err = u.Repo.InTransaction(func(txRepo userr.Repository) error {
		// Проверить, не связан ли пользователь провайдера с каким-нибудь нашим пользователем
		if conflictUsers, err := u.Repo.List(userr.Filter{
			OauthUserID:   openAuthUser.ID,
			OauthProvider: openAuthUser.Provider,
		}); err != nil {
			return err // Возвращает ошибку, если запрос к репозиторию не удался
		} else if len(conflictUsers) != 0 {
			return ErrProvidersUserIsAlreadyLinked // Возвращает ошибку, если пользователь уже связан
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

func (u *OauthCompleteUsecase) login(user userr.User) (Out, error) {
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
