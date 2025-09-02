package chatt

import (
	"time"

	"github.com/google/uuid"
)

// EventInvitationRemoved описывает событие удаления приглашения
type EventInvitationRemoved struct {
	CreatedIn  time.Time
	Recipients []uuid.UUID
	Invitation Invitation
}

// EventInvitationAdded описывает событие добавления приглашения
type EventInvitationAdded struct {
	CreatedIn  time.Time
	Recipients []uuid.UUID
	Invitation Invitation
}

// EventParticipantAdded описывает событие добавления участника
type EventParticipantAdded struct {
	CreatedIn   time.Time
	Recipients  []uuid.UUID
	ChatID      uuid.UUID
	Participant Participant
}

// EventParticipantRemoved описывает событие удаления участника
type EventParticipantRemoved struct {
	CreatedIn   time.Time
	Recipients  []uuid.UUID
	ChatID      uuid.UUID
	Participant Participant
}

// EventChatNameUpdated описывает событие обновления названия чата
type EventChatNameUpdated struct {
	CreatedIn  time.Time
	Recipients []uuid.UUID
	ChatID     uuid.UUID
	Name       string
}
