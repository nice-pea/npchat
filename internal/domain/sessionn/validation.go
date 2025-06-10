package sessionn

import (
	"golang.org/x/exp/slices"
)

var allSessionStatuses = []string{
	StatusNew,
	StatusVerified,
	StatusExpired,
	StatusRevoked,
}

// ValidateSessionStatus проверяет корректность статуса сессии.
func ValidateSessionStatus(status string) error {
	if !slices.Contains(allSessionStatuses, status) {
		return ErrSessionStatusValidate
	}

	return nil
}

// ValidateSessionName проверяет корректность названия сессии.
func ValidateSessionName(name string) error {
	if name == "" {
		return ErrSessionNameEmpty
	}

	return nil
}
