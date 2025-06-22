package domain

import (
	"github.com/google/uuid"
)

// ValidateID валидирует  ID как uuid
func ValidateID(id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrInvalidID
	}

	return nil
}
