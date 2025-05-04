package domain

import (
	"testing"

	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

func TestUser_ValidateID(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(ID string) error {
		u := User{ID: ID}
		return u.ValidateID()
	})
}
