package domain

import (
	"errors"

	"github.com/google/uuid"
)

// Invitation представляет собой отправленное приглашение в чат.
type Invitation struct {
	ID            string // Уникальный идентификатор приглашения
	SubjectUserID string // Пользователь, отправивший приглашение
	UserID        string // Пользователь, получивший приглашение
	ChatID        string // Чата, в который приглашают
}

// Ошибки, связанные с валидацией приглашения.
var (
	ErrInvitationIDValidate            = errors.New("некорректный UUID")
	ErrInvitationChatIDValidate        = errors.New("некорректный ChatID")
	ErrInvitationUserIDValidate        = errors.New("некорректный UserID")
	ErrInvitationSubjectUserIDValidate = errors.New("некорректный SubjectID")
)

// ValidateID проверяет корректность идентификатора приглашения.
func (i Invitation) ValidateID() error {
	if err := uuid.Validate(i.ID); err != nil {
		return errors.Join(err, ErrInvitationIDValidate)
	}
	return nil // Идентификатор валиден
}

// ValidateChatID проверяет корректность идентификатора чата.
func (i Invitation) ValidateChatID() error {
	if err := uuid.Validate(i.ChatID); err != nil {
		return errors.Join(err, ErrInvitationChatIDValidate)
	}
	return nil // Идентификатор чата валиден
}

// ValidateUserID проверяет корректность идентификатора пользователя.
func (i Invitation) ValidateUserID() error {
	if err := uuid.Validate(i.UserID); err != nil {
		return errors.Join(err, ErrInvitationUserIDValidate)
	}
	return nil // Идентификатор пользователя валиден
}

// ValidateSubjectUserID проверяет корректность идентификатора предмета приглашения.
func (i Invitation) ValidateSubjectUserID() error {
	if err := uuid.Validate(i.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvitationSubjectUserIDValidate)
	}
	return nil // Идентификатор предмета приглашения валиден
}

// InvitationsRepository интерфейс для работы с репозиторием приглашений.
type InvitationsRepository interface {
	// List возвращает список с учетом фильтрации
	List(filter InvitationsFilter) ([]Invitation, error)

	// Save сохраняет запись
	Save(invitation Invitation) error

	// Delete удаляет запись
	Delete(id string) error
}

// InvitationsFilter представляет собой фильтр по приглашениям.
type InvitationsFilter struct {
	ID            string // Фильтрация по ID приглашения
	ChatID        string // Фильтрация по чату
	UserID        string // Фильтрация по приглашаемому пользователю
	SubjectUserID string // Фильтрация по пригласившему пользователю
}
