package domain

import (
	"errors"

	"github.com/google/uuid"
)

// Session представляет собой структуру для хранения информации о сессии пользователя.
type Session struct {
	ID     string // ID сессии
	UserID string // ID пользователя, к которому относится сессия
	Token  string // Токен сессии для аутентификации
	Status int    // Статус сессии
}

// Ошибки, связанные с валидацией сессии.
var (
	ErrSessionStatusValidate = errors.New("некорректный статус сессии")
	ErrSessionIDValidate     = errors.New("некорректный ID")
	ErrSessionUserIDValidate = errors.New("некорректный UserID")
)

// ValidateID проверяет корректность идентификатора сессии.
func (s Session) ValidateID() error {
	if err := uuid.Validate(s.ID); err != nil {
		return errors.Join(err, ErrSessionIDValidate)
	}
	return nil // Идентификатор валиден
}

// ValidateUserID проверяет корректность идентификатора пользователя.
func (s Session) ValidateUserID() error {
	if err := uuid.Validate(s.UserID); err != nil {
		return errors.Join(err, ErrSessionUserIDValidate)
	}
	return nil // Идентификатор пользователя валиден
}

// ValidateStatus проверяет корректность статуса сессии.
func (s Session) ValidateStatus() error {
	if s.Status < SessionStatusNew || s.Status > SessionStatusFailed {
		return ErrSessionStatusValidate // Возвращает ошибку, если статус сессии некорректен
	}
	return nil // Статус валиден
}

// Константы, представляющие возможные статусы сессии.
const (
	SessionStatusNew      = 1 // Статус "Новая"
	SessionStatusPending  = 2 // Статус "В ожидании"
	SessionStatusVerified = 3 // Статус "Подтвержденная"
	SessionStatusExpired  = 4 // Статус "Истекшая"
	SessionStatusRevoked  = 5 // Статус "Отозванная"
	SessionStatusFailed   = 6 // Статус "Неудачная"
)

// SessionsRepository интерфейс для работы с репозиторием сессий.
type SessionsRepository interface {
	// Save сохраняет запись
	Save(session Session) error

	// List возвращает список с учетом фильтрации
	List(filter SessionsFilter) ([]Session, error)

	// Delete удаляет запись
	Delete(id string) error
}

// SessionsFilter представляет собой фильтр по сессиям.
type SessionsFilter struct {
	Token string // Фильтрация по токену сессии
}
