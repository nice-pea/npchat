package chatt

import (
	"time"

	"github.com/google/uuid"
)

type EventInvitationRemoved struct {
	CreatedIn  time.Time
	Recipients []uuid.UUID
	Invitation Invitation
}

type EventParticipantAdded struct {
	CreatedIn   time.Time
	Recipients  []uuid.UUID
	ChatID      uuid.UUID
	Participant Participant
}
