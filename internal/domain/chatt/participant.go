package chatt

import (
	"errors"
	"time"

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
func (c *Chat) RemoveParticipant(userID uuid.UUID) error {
	// Убедиться, что участник не является главным администратором
	if userID == c.ChiefID {
		return ErrCannotRemoveChief
	}

	// Убедиться, что участник существует
	if !c.HasParticipant(userID) {
		return ErrParticipantNotExists
	}

	// Удалить участника из списка
	c.Participants = slices.DeleteFunc(c.Participants, func(p Participant) bool {
		return p.UserID == userID
	})

	return nil
}

// AddParticipant добавляет участника в чат.
func (c *Chat) AddParticipant(p Participant, events *events.Events) error {
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
	events.AddSafety(EventParticipantAdded{
		CreatedIn:   time.Now(),
		Recipients:  userIDs(c.Participants),
		ChatID:      c.ID,
		Participant: p,
	})

	return nil
}
