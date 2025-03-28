package domain

import (
	"testing"
)

func TestMember_ValidateID(t *testing.T) {
	RunValidateRequiredIDTest(t, func(ID string) error {
		m := Member{ID: ID}
		return m.ValidateID()
	})
}

func TestMember_ValidateChatID(t *testing.T) {
	RunValidateRequiredIDTest(t, func(ChatID string) error {
		m := Member{ChatID: ChatID}
		return m.ValidateChatID()
	})
}

func TestMember_ValidateUserID(t *testing.T) {
	RunValidateRequiredIDTest(t, func(UserID string) error {
		m := Member{UserID: UserID}
		return m.ValidateUserID()
	})
}
