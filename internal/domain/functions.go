package domain

import (
	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/common"
)

// ValidateID валидирует  ID как uuid
func ValidateID(id uuid.UUID) error {
	if common.IsZero(id) {
		return ErrInvalidID
	}

	return nil
}
