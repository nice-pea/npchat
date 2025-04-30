package domain

import (
	"errors"

	"github.com/google/uuid"
)

type Session struct {
	ID     string
	UserID string
	Token  string
	Status int
}

var (
	ErrSessionStatusValidate = errors.New("некорректный статус сессии")
	ErrSessionIDValidate     = errors.New("некорректный ID")
	ErrSessionUserIDValidate = errors.New("некорректный UserID")
)

func (s Session) ValidateID() error {
	if err := uuid.Validate(s.ID); err != nil {
		return errors.Join(err, ErrSessionIDValidate)
	}
	return nil
}

func (s Session) ValidateUserID() error {
	if err := uuid.Validate(s.UserID); err != nil {
		return errors.Join(err, ErrSessionUserIDValidate)
	}
	return nil
}

func (s Session) ValidateStatus() error {
	if s.Status < SessionStatusNew || s.Status > SessionStatusFailed {
		return ErrSessionStatusValidate
	}
	return nil
}

const (
	SessionStatusNew      = 1
	SessionStatusPending  = 2
	SessionStatusVerified = 3
	SessionStatusExpired  = 4
	SessionStatusRevoked  = 5
	SessionStatusFailed   = 6
)

type SessionsRepository interface {
	Save(session Session) error
	List(filter SessionFilter) ([]Session, error)
	Delete(id string) error
}

type SessionFilter struct {
	Token string
}
