package domain

import (
	"errors"

	"github.com/google/uuid"
)

// User представляет собой пользователя.
type User struct {
	ID   string // ID пользователя
	Name string // Имя пользователя
	Nick string // Ник пользователя
}

var (
	ErrUserIDValidate = errors.New("некорректный UUID")
)

// ValidateID проверяет корректность идентификатора пользователя.
func (u User) ValidateID() error {
	if err := uuid.Validate(u.ID); err != nil {
		return errors.Join(err, ErrUserIDValidate)
	}
	return nil // Идентификатор валиден
}

// UsersRepository интерфейс для работы с репозиторием пользователей.
type UsersRepository interface {
	// List возвращает список с учетом фильтрации
	List(filter UsersFilter) ([]User, error)

	// Save сохраняет запись
	Save(user User) error

	// Delete удаляет запись
	Delete(id string) error
}

// UsersFilter представляет собой фильтр по пользователям.
type UsersFilter struct {
	ID string // ID пользователя для фильтрации
}
