package chatt

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Invitation представляет собой отправленное приглашение в чат.
type Invitation struct {
	ID          string // Глобальный уникальный ID приглашения
	RecipientID string // Пользователь, получивший приглашение
	SubjectID   string // Пользователь, отправивший приглашение
}

// NewInvitation создает новое приглашение в чате
func NewInvitation(subjectID, recipientID string) (Invitation, error) {
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
		ID:          uuid.NewString(),
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
		return ErrUserIsAlreadyInChat
	}

	// Проверить, не существует ли приглашение для этого пользователя в этот чат
	if c.HasInvitationWithRecipient(invitation.RecipientID) {
		return ErrUserIsAlreadyInvited
	}

	c.Invitations = append(c.Invitations, invitation)

	return nil
}

// RemoveInvitation удаляет приглашение из чата
func (c *Chat) RemoveInvitation(id string) error {
	// Убедиться, что приглашение существует
	if !c.HasInvitation(id) {
		return ErrInvitationNotExists
	}

	// Удалить приглашение из списка
	c.Invitations = slices.DeleteFunc(c.Invitations, func(i Invitation) bool {
		return i.ID == id
	})

	return nil
}

// Invitation возвращает приглашение по его ID
func (c *Chat) Invitation(id string) (Invitation, error) {
	for _, i := range c.Invitations {
		if i.ID == id {
			return i, nil
		}
	}

	return Invitation{}, ErrInvitationNotExists
}

// HasInvitation проверяет, существует ли приглашение с указанным ID
func (c *Chat) HasInvitation(id string) bool {
	_, err := c.Invitation(id)

	// Если приглашение найдено, возвращаем true
	return err == nil
}

// HasInvitationWithRecipient проверяет, существует ли приглашение с указанным получателем
func (c *Chat) HasInvitationWithRecipient(recipientID string) bool {
	for _, i := range c.Invitations {
		if i.RecipientID == recipientID {
			return true
		}
	}

	return false
}
