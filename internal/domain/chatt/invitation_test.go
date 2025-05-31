package chatt

import (
	"testing"

	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

func TestInvitation_ValidateID(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(ID string) error {
		i := Invitation{ID: ID}
		return i.ValidateID()
	})
}

func TestInvitation_ValidateChatID(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(ChatID string) error {
		i := Invitation{ChatID: ChatID}
		return i.ValidateChatID()
	})
}

func TestInvitation_ValidateUserID(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(UserID string) error {
		i := Invitation{UserID: UserID}
		return i.ValidateUserID()
	})
}

func TestInvitation_ValidateSubjectUserID(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(SubjectUserID string) error {
		i := Invitation{SubjectUserID: SubjectUserID}
		return i.ValidateSubjectUserID()
	})
}
