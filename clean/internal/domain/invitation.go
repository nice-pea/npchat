package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Invitation struct {
	ID string
	// UserID string
	ChatID string `db:"chat_id"`
}

var (
	ErrInvitationIDValidate     = errors.New("некорректный UUID")
	ErrInvitationChatIDValidate = errors.New("некорректный ChatID")
)

func (i Invitation) ValidateID() error {
	if err := uuid.Validate(i.ID); err != nil {
		return errors.Join(err, ErrInvitationIDValidate)
	}

	return nil
}

func (i Invitation) ValidateChatID() error {
	panic("unimplemented")
}

type InvitationsRepository interface {
	List(filter InvitationsFilter) ([]Invitation, error)
	Save(invitation Invitation) error
	Delete(id string) error
}

type InvitationsFilter struct {
	ID     string
	ChatID string
}
