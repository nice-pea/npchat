package userr

import "time"

// OpenAuthLink представляет собой связь между нашим пользователем и пользователем OAuth провайдера.
type OpenAuthLink struct {
	ExternalID string // ID пользователя провайдером
	Provider   string // Провайдер, которому принадлежит пользователь
	UserID     string // ID нашего пользователя
	Token      OpenAuthToken
}

// OpenAuthToken представляет собой OAuth токен.
type OpenAuthToken struct {
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
}
