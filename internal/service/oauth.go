package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// OAuth сервис, объединяющий случаи использования(юзкейсы) в контексте сущности
type OAuth struct {
	Providers    OAuthProviders            // Карта провайдеров OAuth
	OAuthRepo    domain.OAuthRepository    // Репозиторий для OAuth токенов
	UsersRepo    domain.UsersRepository    // Репозиторий пользователей
	SessionsRepo domain.SessionsRepository // Репозиторий сессий пользователей
}

// OAuthProviders представляет собой карту провайдеров OAuth, где ключом является имя провайдера.
type OAuthProviders = map[string]OAuthProvider

// OAuthProvider определяет интерфейс для работы с провайдерами OAuth.
type OAuthProvider interface {
	// Exchange обменивает код авторизации на токен OAuth
	Exchange(code string) (domain.OAuthToken, error)

	// User возвращает информацию о пользователе провайдера, используя токен OAuth
	User(token domain.OAuthToken) (domain.OAuthUser, error)

	// AuthorizationURL возвращает URL для авторизации.
	// Параметр state используется для предотвращения CSRF-атаки, Должен быть уникальной случайной строкой
	AuthorizationURL(state string) string

	// Name возвращает имя провайдера OAuth
	Name() string
}

// OAuthCompeteRegistrationInput представляет собой структуру для входных данных завершения регистрации OAuth.
type OAuthCompeteRegistrationInput struct {
	UserCode string // Код пользователя, полученный от провайдера
	Provider string // Имя провайдера OAuth
}

// Validate валидирует значение каждого параметра.
func (in OAuthCompeteRegistrationInput) Validate() error {
	if in.UserCode == "" {
		return ErrInvalidUserCode
	}
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	return nil // Возвращает nil, если все параметры валидны
}

// OAuthCompeteRegistrationOut представляет собой результат завершения регистрации OAuth.
type OAuthCompeteRegistrationOut struct {
	Session domain.Session // Сессия пользователя
	User    domain.User    // Пользователь
}

// CompeteRegistration завершает процесс регистрации пользователя через OAuth.
func (o *OAuth) CompeteRegistration(in OAuthCompeteRegistrationInput) (OAuthCompeteRegistrationOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Определить провайдера OAuth
	provider, err := o.provider(in.Provider)
	if err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Получить пользователя провайдера
	token, err := provider.Exchange(in.UserCode)
	if err != nil {
		return OAuthCompeteRegistrationOut{}, errors.Join(ErrWrongUserCode, err)
	}
	pUser, err := provider.User(token)
	if err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Проверить, не связан ли пользователь провайдера с каким-нибудь нашим пользователем
	links, err := o.OAuthRepo.ListLinks(domain.OAuthListLinksFilter{
		ExternalID: pUser.ID,
		Provider:   in.Provider,
	})
	if err != nil {
		return OAuthCompeteRegistrationOut{}, err // Возвращает ошибку, если запрос к репозиторию не удался
	}
	if len(links) != 0 {
		return OAuthCompeteRegistrationOut{}, ErrProvidersUserIsAlreadyLinked // Возвращает ошибку, если пользователь уже связан
	}

	// Создать пользователя
	user := domain.User{
		ID:   uuid.NewString(),
		Name: pUser.Name,
		Nick: "",
	}
	if err = o.UsersRepo.Save(user); err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Сохранить связь нашего пользователя с пользователем провайдера
	err = o.OAuthRepo.SaveLink(domain.OAuthLink{
		UserID:     user.ID,
		ExternalID: pUser.ID,
		Provider:   in.Provider, // Имя провайдера
	})
	if err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	// Создать сессию для пользователя
	session := domain.Session{
		ID:     uuid.NewString(),
		UserID: user.ID,
		Token:  uuid.NewString(),
		Status: domain.SessionStatusVerified, // Подтвержденная сессия
	}
	if err = o.SessionsRepo.Save(session); err != nil {
		return OAuthCompeteRegistrationOut{}, err
	}

	return OAuthCompeteRegistrationOut{
		Session: session,
		User:    user,
	}, nil
}

// OAuthRegistrationInitOut представляет собой результат инициализации регистрации OAuth.
type OAuthRegistrationInitOut struct {
	RedirectURL string // URL для перенаправления на страницу авторизации провайдера
}

// OAuthInitRegistrationInput представляет собой параметры инициализации регистрации OAuth.
type OAuthInitRegistrationInput struct {
	Provider string // Имя провайдера OAuth
}

// Validate валидирует значение параметра провайдера.
func (in OAuthInitRegistrationInput) Validate() error {
	if in.Provider == "" {
		return ErrInvalidProvider
	}

	return nil // Возвращает nil, если параметры валидны
}

// InitRegistration инициализирует процесс регистрации пользователя через OAuth.
func (o *OAuth) InitRegistration(in OAuthInitRegistrationInput) (OAuthRegistrationInitOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return OAuthRegistrationInitOut{}, err
	}

	// Определить провайдера OAuth
	provider, err := o.provider(in.Provider)
	if err != nil {
		return OAuthRegistrationInitOut{}, err
	}

	// Генерирует URL для перенаправления на страницу авторизации провайдера
	return OAuthRegistrationInitOut{
		RedirectURL: provider.AuthorizationURL(uuid.NewString()),
	}, nil
}

// provider возвращает провайдера OAuth по его имени.
func (o *OAuth) provider(provider string) (OAuthProvider, error) {
	p, ok := o.Providers[provider]
	// Проверить, существует ли провайдер в карте
	if !ok || p == nil {
		return nil, ErrUnknownOAuthProvider
	}

	return p, nil
}

// OAuthInitLoginInput представляет собой параметры инициализации входа через OAuth.
type OAuthInitLoginInput struct {
	Provider string // Имя провайдера OAuth
}

// OAuthInitLoginOut представляет собой результат инициализации входа через OAuth.
type OAuthInitLoginOut struct {
	RedirectURL string // URL для перенаправления на страницу авторизации провайдера
}

// InitLogin инициализирует процесс входа пользователя через OAuth.
func (o *OAuth) InitLogin(in OAuthInitLoginInput) (OAuthInitLoginOut, error) {
	// TODO
	return OAuthInitLoginOut{}, nil
}

// OAuthCompleteLoginInput представляет собой параметры завершения входа через OAuth.
type OAuthCompleteLoginInput struct {
	UserCode string // Код пользователя, полученный от провайдера
	Provider string // Имя провайдера OAuth
}

// OAuthCompleteLoginOut представляет собой результат завершения входа через OAuth.
type OAuthCompleteLoginOut struct {
	Session domain.Session // Сессия пользователя
	User    domain.User    // Пользователь
}

// CompleteLogin завершает процесс входа пользователя через OAuth.
func (o *OAuth) CompleteLogin(in OAuthCompleteLoginInput) (OAuthCompleteLoginOut, error) {
	// TODO
	return OAuthCompleteLoginOut{}, nil
}
