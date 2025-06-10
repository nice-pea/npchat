package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
)

// OAuthProviders представляет собой карту провайдеров OAuth, где ключом является имя провайдера.
type OAuthProviders = map[string]OAuthProvider

// OAuthProvider определяет интерфейс для работы с провайдерами OAuth.
type OAuthProvider interface {
	// Exchange обменивает код авторизации на токен OAuth
	Exchange(code string) (userr.OpenAuthToken, error)

	// User возвращает информацию о пользователе провайдера, используя токен OAuth
	User(token userr.OpenAuthToken) (userr.OpenAuthUser, error)

	// AuthorizationURL возвращает URL для авторизации.
	// Параметр state используется для предотвращения CSRF-атаки, Должен быть уникальной случайной строкой
	AuthorizationURL(state string) string

	// Name возвращает имя провайдера OAuth
	Name() string
}

// CompeteOAuthRegistrationIn представляет собой структуру для входных данных завершения регистрации OAuth.
type CompeteOAuthRegistrationIn struct {
	UserCode string // Код пользователя, полученный от провайдера
	Provider string // Имя провайдера OAuth
}

// Validate валидирует значение каждого параметра.
func (in CompeteOAuthRegistrationIn) Validate() error {
	if in.UserCode == "" {
		return ErrInvalidUserCode
	}
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	return nil // Возвращает nil, если все параметры валидны
}

// CompeteOAuthRegistrationOut представляет собой результат завершения регистрации OAuth.
type CompeteOAuthRegistrationOut struct {
	Session sessionn.Session // Сессия пользователя
	User    userr.User       // Пользователь
}

// CompeteOAuthRegistration завершает процесс регистрации пользователя через OAuth.
func (u *Users) CompeteOAuthRegistration(in CompeteOAuthRegistrationIn) (CompeteOAuthRegistrationOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return CompeteOAuthRegistrationOut{}, err
	}

	// Определить провайдера OAuth
	provider, err := u.provider(in.Provider)
	if err != nil {
		return CompeteOAuthRegistrationOut{}, err
	}
	// Получить пользователя провайдера
	token, err := provider.Exchange(in.UserCode)
	if err != nil {
		return CompeteOAuthRegistrationOut{}, errors.Join(ErrWrongUserCode, err)
	}
	openAuthUser, err := provider.User(token)
	if err != nil {
		return CompeteOAuthRegistrationOut{}, err
	}

	// Создать пользователя
	user, err := userr.NewUser(openAuthUser.Name, "")
	if err != nil {
		return CompeteOAuthRegistrationOut{}, err
	}
	// Добавить пользователя в список связей для аутентификации по OAuth
	if err = user.AddOpenAuthUser(openAuthUser); err != nil {
		return CompeteOAuthRegistrationOut{}, err
	}

	// Проверить, не связан ли пользователь провайдера с каким-нибудь нашим пользователем
	if conflictUsers, err := u.Repo.List(userr.Filter{
		OAuthUserID:   openAuthUser.ID,
		OAuthProvider: in.Provider,
	}); err != nil {
		return CompeteOAuthRegistrationOut{}, err // Возвращает ошибку, если запрос к репозиторию не удался
	} else if len(conflictUsers) != 0 {
		return CompeteOAuthRegistrationOut{}, ErrProvidersUserIsAlreadyLinked // Возвращает ошибку, если пользователь уже связан
	}

	// Сохранить пользователя в репозиторий
	if err = u.Repo.Upsert(user); err != nil {
		return CompeteOAuthRegistrationOut{}, err
	}

	// Создать сессию для пользователя
	sessionName := "todo: [название модели телефона / название браузера]"
	session, err := sessionn.NewSession(user.ID, sessionName, sessionn.StatusVerified)
	if err != nil {
		return CompeteOAuthRegistrationOut{}, err
	}
	if err = u.SessionsRepo.Upsert(session); err != nil {
		return CompeteOAuthRegistrationOut{}, err
	}

	return CompeteOAuthRegistrationOut{
		Session: session,
		User:    user,
	}, nil
}

// InitOAuthRegistrationOut представляет собой результат инициализации регистрации OAuth.
type InitOAuthRegistrationOut struct {
	RedirectURL string // URL для перенаправления на страницу авторизации провайдера
}

// InitOAuthRegistrationIn представляет собой параметры инициализации регистрации OAuth.
type InitOAuthRegistrationIn struct {
	Provider string // Имя провайдера OAuth
}

// Validate валидирует значение параметра провайдера.
func (in InitOAuthRegistrationIn) Validate() error {
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	return nil // Возвращает nil, если параметры валидны
}

// InitOAuthRegistration инициализирует процесс регистрации пользователя через OAuth.
func (u *Users) InitOAuthRegistration(in InitOAuthRegistrationIn) (InitOAuthRegistrationOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return InitOAuthRegistrationOut{}, err
	}

	// Определить провайдера OAuth
	provider, err := u.provider(in.Provider)
	if err != nil {
		return InitOAuthRegistrationOut{}, err
	}

	// Генерирует URL для перенаправления на страницу авторизации провайдера
	return InitOAuthRegistrationOut{
		RedirectURL: provider.AuthorizationURL(uuid.NewString()),
	}, nil
}

// provider возвращает провайдера OAuth по его имени.
func (u *Users) provider(provider string) (OAuthProvider, error) {
	p, ok := u.Providers[provider]
	// Проверить, существует ли провайдер в карте
	if !ok || p == nil {
		return nil, ErrUnknownOAuthProvider
	}

	return p, nil
}

// InitOAuthLoginIn представляет собой параметры инициализации входа через OAuth.
type InitOAuthLoginIn struct {
	Provider string // Имя провайдера OAuth
}

// InitOAuthLoginOut представляет собой результат инициализации входа через OAuth.
type InitOAuthLoginOut struct {
	RedirectURL string // URL для перенаправления на страницу авторизации провайдера
}

// InitOAuthLogin инициализирует процесс входа пользователя через OAuth.
func (u *Users) InitOAuthLogin(in InitOAuthLoginIn) (InitOAuthLoginOut, error) {
	// TODO
	return InitOAuthLoginOut{}, nil
}

// CompleteOAuthLoginIn представляет собой параметры завершения входа через OAuth.
type CompleteOAuthLoginIn struct {
	UserCode string // Код пользователя, полученный от провайдера
	Provider string // Имя провайдера OAuth
}

// CompleteOAuthLoginOut представляет собой результат завершения входа через OAuth.
type CompleteOAuthLoginOut struct {
	Session sessionn.Session // Сессия пользователя
	User    userr.User       // Пользователь
}

// CompleteOAuthLogin завершает процесс входа пользователя через OAuth.
func (u *Users) CompleteOAuthLogin(in CompleteOAuthLoginIn) (CompleteOAuthLoginOut, error) {
	// TODO
	return CompleteOAuthLoginOut{}, nil
}
