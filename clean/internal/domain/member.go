package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Member struct {
	ID string
	//UserID string
	ChatID string
}

var (
	ErrMemberIDValidate = errors.New("некорректный UUID")
)

func (m Member) ValidateID() error {
	if err := uuid.Validate(m.ID); err != nil {
		return errors.Join(err, ErrMemberIDValidate)
	}

	return nil
}

func (m Member) ValidateChatID() error {
	if err := uuid.Validate(m.ChatID); err != nil {
		return err
	}

	return nil
}

type MembersRepository interface {
	List(filter MembersFilter) ([]Member, error)
	Save(member Member) error
	Delete(id string) error
}

type MembersFilter struct {
	ID string
	//UserID string
	ChatID string
	//IsOwner *bool
}
