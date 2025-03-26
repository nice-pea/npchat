package domain

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

type Chat struct {
	ID          string
	Name        string
	ChiefUserID string
}

var (
	ErrChatIDValidate          = errors.New("некорректный UUID")
	ErrChatChiefUserIDValidate = errors.New("некорректный ChiefUserID")
	ErrChatNameValidate        = errors.New("некорректный Name")
)

func (c Chat) ValidateID() error {
	if err := uuid.Validate(c.ID); err != nil {
		return errors.Join(err, ErrChatIDValidate)
	}

	return nil
}

func (c Chat) ValidateName() error {
	var chatNameRegexp = regexp.MustCompile(`^[^\s\n\t][^\n\t]{0,48}[^\s\n\t]$`)
	if !chatNameRegexp.MatchString(c.Name) {
		return ErrChatNameValidate
	}

	return nil
}

func (c Chat) ValidateChiefUserID() error {
	if err := uuid.Validate(c.ChiefUserID); err != nil {
		return errors.Join(err, ErrChatChiefUserIDValidate)
	}

	return nil
}

type ChatsRepository interface {
	List(filter ChatsFilter) ([]Chat, error)
	Save(chat Chat) error
	Delete(id string) error
}

type ChatsFilter struct {
	ID string
	//UserIDs []string
}
