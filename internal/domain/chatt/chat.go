package chatt

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/usecases/events"
)

// Chat представляет собой агрегат чата.
type Chat struct {
	ID      uuid.UUID // Уникальный ID чата
	Name    string    // Название чата
	ChiefID uuid.UUID // ID главного пользователя чата

	Participants []Participant // Список участников чата
	Invitations  []Invitation  // Список приглашений в чате
}

// NewChat создает новый чат.
func NewChat(name string, chiefID uuid.UUID, eventsBuf *events.Buffer) (Chat, error) {
	if err := ValidateChatName(name); err != nil {
		return Chat{}, err
	}
	if err := domain.ValidateID(chiefID); err != nil {
		return Chat{}, errors.Join(err, ErrInvalidChiefID)
	}

	chat := Chat{
		ID:      uuid.New(),
		Name:    name,
		ChiefID: chiefID,
		Participants: []Participant{
			{UserID: chiefID}, // Главный администратор
		},
		Invitations: nil,
	}

	eventsBuf.AddSafety(EventChatCreated{
		Head:   events.NewHead(userIDs(chat.Participants)),
		ChatID: chat.ID,
	})

	return chat, nil
}

// UpdateName изменяет название чата.
func (c *Chat) UpdateName(name string, eventsBuf *events.Buffer) error {
	if err := ValidateChatName(name); err != nil {
		return err
	}

	c.Name = name

	eventsBuf.AddSafety(EventChatNameUpdated{
		Head:   events.NewHead(userIDs(c.Participants)),
		ChatID: c.ID,
		Name:   c.Name,
	})

	return nil
}

func userIDs(participants []Participant) []uuid.UUID {
	userIDs := make([]uuid.UUID, len(participants))
	for i, p := range participants {
		userIDs[i] = p.UserID
	}
	return userIDs
}
