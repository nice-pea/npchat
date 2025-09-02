package chatt

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/usecases/events"
)

// Participant представляет собой участника чата.
type Participant struct {
	UserID uuid.UUID // ID пользователя, который является участником чата
}

// NewParticipant создает новый участник чата.
func NewParticipant(userID uuid.UUID) (Participant, error) {
	if err := domain.ValidateID(userID); err != nil {
		return Participant{}, errors.Join(err, ErrInvalidUserID)
	}

	return Participant{
		UserID: userID,
	}, nil
}

// HasParticipant проверяет, является ли пользователь участником чата.
func (c *Chat) HasParticipant(userID uuid.UUID) bool {
	for _, p := range c.Participants {
		if p.UserID == userID {
			return true
		}
	}

	return false
}

// RemoveParticipant удаляет участника из чата.
func (c *Chat) RemoveParticipant(userID uuid.UUID, eventsBuf *events.Buffer) error {
	// Убедиться, что участник не является главным администратором
	if userID == c.ChiefID {
		return ErrCannotRemoveChief
	}

	// Убедиться, что участник существует
	if !c.HasParticipant(userID) {
		return ErrParticipantNotExists
	}

	// Найти индекс участника
	i := slices.IndexFunc(c.Participants, func(p Participant) bool {
		return p.UserID == userID
	})

	// Добавить событие
	eventsBuf.AddSafety(EventParticipantRemoved{
		Head:        events.NewHead(userIDs(c.Participants)),
		ChatID:      c.ID,
		Participant: c.Participants[i],
	})

	// Удалить участника
	c.Participants = slices.Delete(c.Participants, i, i+1)

	return nil
}

// AddParticipant добавляет участника в чат.
func (c *Chat) AddParticipant(p Participant, eventsBuf *events.Buffer) error {
	// Проверить является ли subject участником чата
	if c.HasParticipant(p.UserID) {
		return ErrParticipantExists
	}

	// Проверить, не существует ли приглашение для этого пользователя в этот чат
	if c.HasInvitationWithRecipient(p.UserID) {
		return ErrUserIsAlreadyInvited
	}

	// Добавить участника
	c.Participants = append(c.Participants, p)

	// Добавить событие
	eventsBuf.AddSafety(EventParticipantAdded{
		Head:        events.NewHead(userIDs(c.Participants)),
		ChatID:      c.ID,
		Participant: p,
	})

	return nil
}
