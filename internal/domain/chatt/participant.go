package chatt

import (
	"golang.org/x/exp/slices"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Participant представляет собой участника чата.
type Participant struct {
	UserID string // ID пользователя, который является участником чата
}

// NewParticipant создает новый участник чата.
func NewParticipant(userID string) (Participant, error) {
	if err := domain.ValidateID(userID); err != nil {
		return Participant{}, err
	}

	return Participant{
		UserID: userID,
	}, nil
}

// HasParticipant проверяет, является ли пользователь участником чата.
func (c *Chat) HasParticipant(userID string) bool {
	for _, p := range c.Participants {
		if p.UserID == userID {
			return true
		}
	}

	return false
}

// RemoveParticipant удаляет участника из чата.
func (c *Chat) RemoveParticipant(userID string) error {
	// Убедиться, что участник не является главным администратором
	if userID == c.ChiefID {
		return ErrSubjectUserShouldNotBeChief
	}

	// Убедиться, что участник существует
	if !c.HasParticipant(userID) {
		return ErrUserIsNotMember
	}

	// Удалить участника из списка
	c.Participants = slices.DeleteFunc(c.Participants, func(p Participant) bool {
		return p.UserID == userID
	})

	return nil
}

// AddParticipant добавляет участника в чат.
func (c *Chat) AddParticipant(p Participant) error {
	// Проверить является ли subject участником чата
	if c.HasParticipant(p.UserID) {
		return ErrUserIsAlreadyInChat
	}

	// Проверить, не существует ли приглашение для этого пользователя в этот чат
	if c.HasInvitationWithRecipient(p.UserID) {
		return ErrUserIsAlreadyInvited
	}

	c.Participants = append(c.Participants, p)

	return nil
}
