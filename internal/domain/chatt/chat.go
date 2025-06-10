package chatt

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Chat представляет собой агрегат чата.
type Chat struct {
	ID      string // Уникальный ID чата
	Name    string // Название чата
	ChiefID string // ID главного пользователя чата

	Participants []Participant // Список участников чата
	Invitations  []Invitation  // Список приглашений в чате
}

// NewChat создает новый чат.
func NewChat(name string, chiefID string) (Chat, error) {
	if err := ValidateChatName(name); err != nil {
		return Chat{}, err
	}
	if err := domain.ValidateID(chiefID); err != nil {
		return Chat{}, errors.Join(err, ErrInvalidChiefID)
	}

	return Chat{
		ID:      uuid.NewString(),
		Name:    name,
		ChiefID: chiefID,
		Participants: []Participant{
			{UserID: chiefID}, // Главный администратор
		},
		Invitations: nil,
	}, nil
}

// UpdateName изменяет название чата.
func (c *Chat) UpdateName(name string) error {
	if err := ValidateChatName(name); err != nil {
		return err
	}

	c.Name = name

	return nil
}
