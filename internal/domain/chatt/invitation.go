package chatt

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/usecases/events"
)

// Invitation представляет собой отправленное приглашение в чат.
type Invitation struct {
	ID          uuid.UUID // Глобальный уникальный ID приглашения
	RecipientID uuid.UUID // Пользователь, получивший приглашение
	SubjectID   uuid.UUID // Пользователь, отправивший приглашение
}

// NewInvitation создает новое приглашение в чате
func NewInvitation(subjectID, recipientID uuid.UUID) (Invitation, error) {
	if err := domain.ValidateID(subjectID); err != nil {
		return Invitation{}, err
	}
	if err := domain.ValidateID(recipientID); err != nil {
		return Invitation{}, err
	}

	// Пригласивший и приглашаемый не могут быть одним пользователем
	if recipientID == subjectID {
		return Invitation{}, ErrSubjectAndRecipientMustBeDifferent
	}

	return Invitation{
		ID:          uuid.New(),
		RecipientID: recipientID,
		SubjectID:   subjectID,
	}, nil
}

// AddInvitation добавляет приглашение в чат
func (c *Chat) AddInvitation(invitation Invitation) error {
	// Проверить является ли subject участником чата
	if !c.HasParticipant(invitation.SubjectID) {
		return ErrSubjectIsNotMember
	}

	// Проверить является ли user участником чата
	if c.HasParticipant(invitation.RecipientID) {
		return ErrParticipantExists
	}

	// Проверить, не существует ли приглашение для этого пользователя в этот чат
	if c.HasInvitationWithRecipient(invitation.RecipientID) {
		return ErrUserIsAlreadyInvited
	}

	c.Invitations = append(c.Invitations, invitation)

	return nil
}

// RemoveInvitation удаляет приглашение из чата
func (c *Chat) RemoveInvitation(id uuid.UUID, events *events.Buffer) error {
	// Убедиться, что приглашение существует
	if !c.HasInvitation(id) {
		return ErrInvitationNotExists
	}

	// Найти индекс приглашения
	i := slices.IndexFunc(c.Invitations, func(i Invitation) bool {
		return i.ID == id
	})

	// Добавить событие
	events.AddSafety(EventInvitationRemoved{
		CreatedIn: time.Now(),
		Recipients: []uuid.UUID{
			c.ChiefID,
			c.Invitations[i].RecipientID,
			c.Invitations[i].SubjectID,
		},
		Invitation: c.Invitations[i],
	})

	// Удалить приглашение из списка
	c.Invitations = slices.Delete(c.Invitations, i, i+1)

	return nil
}

// Invitation возвращает приглашение по его ID
func (c *Chat) Invitation(id uuid.UUID) (Invitation, error) {
	for _, i := range c.Invitations {
		if i.ID == id {
			return i, nil
		}
	}

	return Invitation{}, ErrInvitationNotExists
}

// HasInvitation проверяет, существует ли приглашение с указанным ID
func (c *Chat) HasInvitation(id uuid.UUID) bool {
	_, err := c.Invitation(id)

	// Если приглашение найдено, возвращаем true
	return err == nil
}

// HasInvitationWithRecipient проверяет, существует ли приглашение с указанным получателем
func (c *Chat) HasInvitationWithRecipient(recipientID uuid.UUID) bool {
	for _, i := range c.Invitations {
		if i.RecipientID == recipientID {
			return true
		}
	}

	return false
}

// SubjectInvitations возвращает список приглашений, отправленных пользователем с указанным ID
func (c *Chat) SubjectInvitations(subjectID uuid.UUID) []Invitation {
	return slices.DeleteFunc(c.Invitations, func(i Invitation) bool {
		return i.SubjectID != subjectID
	})
}

// RecipientInvitation возвращает приглашение, направленное пользователю с указанным ID
func (c *Chat) RecipientInvitation(recipientID uuid.UUID) (Invitation, error) {
	for _, inv := range c.Invitations {
		if inv.RecipientID == recipientID {
			return inv, nil
		}
	}

	return Invitation{}, ErrInvitationNotExists
}
