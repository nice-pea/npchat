package userr

import (
	"errors"
	"fmt"
	"net/url"
	"time"
)

// OpenAuthUser представляет собой пользователем OAuth провайдера.
type OpenAuthUser struct {
	ID       string        // ID пользователя провайдером
	Provider string        // Провайдер, которому принадлежит пользователь
	Email    string        // Электронная почта пользователя
	Name     string        // Имя пользователя
	Picture  string        // URL изображения профиля
	Token    OpenAuthToken // Токен для аутентификации
}

// NewOpenAuthUser создает нового пользователем OAuth провайдера.
func NewOpenAuthUser(id string, provider string, email string, name string, picture string, token OpenAuthToken) (OpenAuthUser, error) {
	if id == "" {
		return OpenAuthUser{}, errors.New("id is required")
	}
	if provider == "" {
		return OpenAuthUser{}, errors.New("provider is required")
	}
	if email != "" {
		if err := ValidateEmail(email); err != nil {
			return OpenAuthUser{}, errors.New("email is required")
		}
	}
	if name == "" {
		return OpenAuthUser{}, errors.New("name is required")
	}
	if picture != "" {
		if _, err := url.Parse(picture); err != nil {
			return OpenAuthUser{}, fmt.Errorf("invalid picture url: %w", err)
		}
	}

	return OpenAuthUser{
		ID:       id,
		Provider: provider,
		Email:    email,
		Name:     name,
		Picture:  picture,
		Token:    token,
	}, nil
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

// NewOpenAuthToken создает новый OAuth токен.
func NewOpenAuthToken(accessToken string, tokenType string, refreshToken string, expiry time.Time) (OpenAuthToken, error) {
	if accessToken == "" {
		return OpenAuthToken{}, errors.New("accessToken is required")
	}
	if tokenType == "" {
		return OpenAuthToken{}, errors.New("tokenType is required")
	}

	return OpenAuthToken{
		AccessToken:  accessToken,
		TokenType:    tokenType,
		RefreshToken: refreshToken,
		Expiry:       expiry,
	}, nil
}
