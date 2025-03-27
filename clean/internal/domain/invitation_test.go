package domain

import (
	"testing"
)

func TestInvitation_ValidateID(t *testing.T) {
	RunValidateRequiredIDTest(t, func(ID string) error {
		i := Invitation{ID: ID}
		return i.ValidateID()
	})
}

func TestInvitation_ValidateChatID(t *testing.T) {
	RunValidateRequiredIDTest(t, func(ChatID string) error {
		i := Invitation{ChatID: ChatID}
		return i.ValidateChatID()
	})
}
