package domain

import (
	"errors"

	"github.com/google/uuid"
)

// Member представляет собой участника чата.
type Member struct {
	ID     string // Глобальный ID участника
	UserID string // ID пользователя, который является участником чата
	ChatID string // ID чата
}

// Ошибки, связанные с валидацией члена чата.
var (
	ErrMemberIDValidate     = errors.New("некорректный UUID")
	ErrMemberChatIDValidate = errors.New("некорректный ChatID")
	ErrMemberUserIDValidate = errors.New("некорректный UserID")
)

// ValidateID проверяет корректность идентификатора члена чата.
func (m Member) ValidateID() error {
	if err := uuid.Validate(m.ID); err != nil {
		return errors.Join(err, ErrMemberIDValidate)
	}
	return nil // Идентификатор валиден
}

// ValidateChatID проверяет корректность идентификатора чата.
func (m Member) ValidateChatID() error {
	if err := uuid.Validate(m.ChatID); err != nil {
		return errors.Join(err, ErrMemberChatIDValidate)
	}
	return nil // Идентификатор чата валиден
}

// ValidateUserID проверяет корректность идентификатора пользователя.
func (m Member) ValidateUserID() error {
	if err := uuid.Validate(m.UserID); err != nil {
		return errors.Join(err, ErrMemberUserIDValidate)
	}
	return nil // Идентификатор пользователя валиден
}

// MembersRepository интерфейс для работы с репозиторием членов чата.
type MembersRepository interface {
	// List возвращает список с учетом фильтрации
	List(filter MembersFilter) ([]Member, error)

	// Save сохраняет запись
	Save(member Member) error

	// Delete удаляет запись
	Delete(id string) error
}

// MembersFilter представляет собой фильтр по членам чата.
type MembersFilter struct {
	ID     string // Фильтрация по ID
	UserID string // Фильтрация по пользователю
	ChatID string // Фильтрация по чату
}
