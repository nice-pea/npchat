package userr

import (
	"errors"
	"time"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// OpenAuthLink представляет собой связь между нашим пользователем и пользователем OAuth провайдера.
type OpenAuthLink struct {
	ExternalID string // ID пользователя провайдером
	Provider   string // Провайдер, которому принадлежит пользователь
	UserID     string // ID нашего пользователя
	Token      OpenAuthToken
}

// NewOpenAuthLink создает новую связь между нашим пользователем и пользователем OAuth провайдера.
func NewOpenAuthLink(externalID string, provider string, userID string, token OpenAuthToken) (OpenAuthLink, error) {
	if externalID == "" {
		return OpenAuthLink{}, errors.New("externalID is required")
	}
	if provider == "" {
		return OpenAuthLink{}, errors.New("provider is required")
	}
	if err := domain.ValidateID(userID); err != nil {
		return OpenAuthLink{}, err
	}

	return OpenAuthLink{
		ExternalID: externalID,
		Provider:   provider,
		UserID:     userID,
		Token:      token,
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
