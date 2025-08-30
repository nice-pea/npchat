package completeOAuthRegistration

import (
	"errors"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	"github.com/nice-pea/npchat/internal/service/users/oauth"
)

var (
	ErrInvalidProvider              = errors.New("некорректное значение Provider")
	ErrProvidersUserIsAlreadyLinked = errors.New("пользователь OAuth-провайдера уже связан с пользователем")
	ErrInvalidUserCode              = errors.New("некорректное значение UserCode")
	ErrWrongUserCode                = errors.New("неправильный UserCode")
)

// In представляет собой структуру для входных данных завершения регистрации OAuth.
type In struct {
	UserCode string // Код пользователя, полученный от провайдера
	Provider string // Имя провайдера OAuth
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

// Out представляет собой результат завершения регистрации OAuth.
type Out struct {
	Session sessionn.Session // Сессия пользователя
	User    userr.User       // Пользователь
}

type CompleteOAuthRegistrationUsecase struct {
	Repo         userr.Repository
	Providers    oauth.OAuthProviders
	SessionsRepo sessionn.Repository // Репозиторий сессий пользователей
}

// CompleteOAuthRegistration завершает процесс регистрации пользователя через OAuth.
func (u *CompleteOAuthRegistrationUsecase) CompleteOAuthRegistration(in In) (Out, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return Out{}, err
	}

	// Определить провайдера OAuth
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

	// Создать пользователя
	user, err := userr.NewUser(openAuthUser.Name, "")
	if err != nil {
		return Out{}, err
	}
	// Добавить пользователя в список связей для аутентификации по OAuth
	if err = user.AddOpenAuthUser(openAuthUser); err != nil {
		return Out{}, err
	}

	// Обернуть работу с репозиторием в транзакцию
	var session sessionn.Session
	if err = u.Repo.InTransaction(func(txRepo userr.Repository) error {
		// Проверить, не связан ли пользователь провайдера с каким-нибудь нашим пользователем
		if conflictUsers, err := u.Repo.List(userr.Filter{
			OAuthUserID:   openAuthUser.ID,
			OAuthProvider: in.Provider,
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
