package domain

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

type Chat struct {
	ID   string
	Name string
}

var (
	ErrChatIDValidate   = errors.New("некорректный UUID. Пожалуйста, введите действительный UUID")
	ErrChatNameValidate = errors.New("название чата должно содержать от 1 до 50 символов, включать только буквы, цифры и пробелы")
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

type ChatsRepository interface {
	List(filter ChatsFilter) ([]Chat, error)
	Save(chat Chat) error
	Delete(id string) error
}

type ChatsFilter struct {
	ID      string
	UserIDs []string
}
