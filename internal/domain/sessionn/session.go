package sessionn

import (
	"time"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Session представляет собой агрегат сессии
type Session struct {
	ID           uuid.UUID // ID сессии
	UserID       uuid.UUID // ID пользователя, к которому относится сессия
	Name         string    // [название модели телефона / название браузера]
	Status       string    // Статус сессии
	AccessToken  Token     // Токен сессии для аутентификации
	RefreshToken Token     // Токен для обновления AccessToken
}

var (
	refreshTokenLifetime = 60 * 24 * time.Hour // 60 дней
	accessTokenLifetime  = 15 * time.Minute    // 15 минут
)

// NewSession создает новую сессию связанную с пользователем.
func NewSession(userID uuid.UUID, name, status string) (Session, error) {
	if err := domain.ValidateID(userID); err != nil {
		return Session{}, err
	}
	if err := ValidateSessionName(name); err != nil {
		return Session{}, err
	}
	if err := ValidateSessionStatus(status); err != nil {
		return Session{}, err
	}

	return Session{
		ID:     uuid.New(),
		UserID: userID,
		Name:   name,
		Status: status,
		AccessToken: Token{
			Token:  uuid.NewString(),
			Expiry: time.Now().Add(accessTokenLifetime),
		},
		RefreshToken: Token{
			Token:  uuid.NewString(),
			Expiry: time.Now().Add(refreshTokenLifetime),
		},
	}, nil
}
