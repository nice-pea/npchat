package domain

import (
	"errors"
	"regexp"

	"github.com/google/uuid"
)

// Chat представляет собой чат.
type Chat struct {
	ID          string // Уникальный идентификатор чата
	Name        string // Название чата
	ChiefUserID string // Идентификатор главного пользователя чата
}

// Ошибки, связанные с валидацией чата.
var (
	ErrChatIDValidate          = errors.New("некорректный UUID")
	ErrChatChiefUserIDValidate = errors.New("некорректный ChiefUserID")
	ErrChatNameValidate        = errors.New("некорректный Name")
)

// ValidateID проверяет корректность идентификатора чата.
func (c Chat) ValidateID() error {
	if err := uuid.Validate(c.ID); err != nil {
		return errors.Join(err, ErrChatIDValidate) // Возвращает ошибку, если идентификатор некорректен
	}

	return nil // Идентификатор валиден
}

// ValidateChatName проверяет корректность названия чата.
func ValidateChatName(name string) error {
	// Регулярное выражение для проверки названия чата
	var chatNameRegexp = regexp.MustCompile(`^[^\s\n\t][^\n\t]{0,48}[^\s\n\t]$`)
	if !chatNameRegexp.MatchString(name) {
		return ErrChatNameValidate // Возвращает ошибку, если название некорректно
	}

	return nil // Название валидно
}

// ChatsRepository интерфейс для работы с репозиторием чатов.
type ChatsRepository interface {
	// List возвращает список с учетом фильтрации
	List(filter ChatsFilter) ([]Chat, error)

	// Save сохраняет запись
	Save(chat Chat) error

	// Delete удаляет запись
	Delete(id string) error
}

// ChatsFilter представляет собой фильтр по чатам.
type ChatsFilter struct {
	IDs []string // Список идентификаторов чатов для фильтрации
}
