package domain

import (
	"errors"

	"github.com/google/uuid"
)

func ValidateID(id string) error {
	if err := uuid.Validate(id); err != nil {
		return errors.Join(err, ErrInvalidID) // Возвращает ошибку, если идентификатор некорректен
	}

	return nil // Идентификатор валиден
}
