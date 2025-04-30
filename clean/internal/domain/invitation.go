package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Invitation struct {
	ID            string
	SubjectUserID string
	UserID        string
	ChatID        string
}

var (
	ErrInvitationIDValidate            = errors.New("некорректный UUID")
	ErrInvitationChatIDValidate        = errors.New("некорректный ChatID")
	ErrInvitationUserIDValidate        = errors.New("некорректный UserID")
	ErrInvitationSubjectUserIDValidate = errors.New("некорректный SubjectUserID")
)

func (i Invitation) ValidateID() error {
	if err := uuid.Validate(i.ID); err != nil {
		return errors.Join(err, ErrInvitationIDValidate)
	}
	return nil
}

func (i Invitation) ValidateChatID() error {
	if err := uuid.Validate(i.ChatID); err != nil {
		return errors.Join(err, ErrInvitationChatIDValidate)
	}
	return nil
}

func (i Invitation) ValidateUserID() error {
	if err := uuid.Validate(i.UserID); err != nil {
		return errors.Join(err, ErrInvitationUserIDValidate)
	}
	return nil
}

func (i Invitation) ValidateSubjectUserID() error {
	if err := uuid.Validate(i.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvitationSubjectUserIDValidate)
	}
	return nil
}

type InvitationsRepository interface {
	List(filter InvitationsFilter) ([]Invitation, error)
	Save(invitation Invitation) error
	Delete(id string) error
}

type InvitationsFilter struct {
	ID            string
	ChatID        string
	UserID        string
	SubjectUserID string
}
