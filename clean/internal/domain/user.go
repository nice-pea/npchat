package domain

import (
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID string
	// Name string
}

var (
	ErrUserIDValidate = errors.New("некорректный UUID")
)

func (u User) ValidateID() error {
	if err := uuid.Validate(u.ID); err != nil {
		return errors.Join(err, ErrUserIDValidate)
	}

	return nil
}

type UsersRepository interface {
	List(filter UsersFilter) ([]User, error)
	Save(user User) error
	Delete(id string) error
}

type UsersFilter struct {
	ID string
}
