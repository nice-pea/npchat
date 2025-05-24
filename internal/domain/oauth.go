package domain

import "time"

// OAuthToken представляет собой OAuth2 токен.
type OAuthToken struct {
	// AccessToken это токен, который авторизует и аутентифицирует запросы.
	AccessToken string

	// TokenType это тип токена.
	TokenType string

	// RefreshToken это токен, который используется приложением
	// (в отличие от пользователя) для обновления токена доступа
	// в случае его истечения.
	RefreshToken string

	// Expiry это необязательное время истечения токена доступа.
	//
	// Если равно нулю, AccessToken будет валидным навсегда,
	// и механизм обновления с помощью RefreshToken не будет использоваться.
	Expiry time.Time

	// LinkID это идентификатор ссылки, связанной с токеном.
	// TODO: Обязательное поле.
	LinkID string

	// Provider это провайдер, выдавший токен.
	Provider string
}

// OAuthUser представляет собой пользователя OAuth2.
type OAuthUser struct {
	ID       string // ID пользователя в провайдере
	Email    string // Электронная почта
	Name     string // Имя пользователя
	Picture  string // URL изображения профиля
	Provider string // Провайдер, которому принадлежит пользователь
}

// OAuthLink представляет собой структуру для хранения информации о связи между пользователем и внешним идентификатором.
type OAuthLink struct {
	UserID     string // ID нашего пользователя
	ExternalID string // ID пользователя провайдером
	Provider   string // Провайдер, которому принадлежит пользователь
}

// OAuthRepository интерфейс для работы с репозиторием OAuth.
type OAuthRepository interface {
	// SaveToken сохраняет запись.
	// TODO: Сейчас не используется; Понадобится для внедрения автообновления токена.
	SaveToken(OAuthToken) error

	// SaveLink сохраняет запись
	SaveLink(OAuthLink) error

	// ListLinks возвращает список с учетом фильтрации
	ListLinks(filter OAuthListLinksFilter) ([]OAuthLink, error)
}

// OAuthListLinksFilter представляет собой фильтр по связям пользователей с провайдерами.
type OAuthListLinksFilter struct {
	UserID     string // Фильтрация по ID пользователя
	ExternalID string // Фильтрация по ID пользователя провайдера
	Provider   string // Фильтрация по провайдеру
}
