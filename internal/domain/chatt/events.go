package chatt

import (
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

// EventInvitationRemoved описывает событие удаления приглашения
type EventInvitationRemoved struct {
	events.Head
	Invitation Invitation
}

// EventInvitationAdded описывает событие добавления приглашения
type EventInvitationAdded struct {
	events.Head
	Invitation Invitation
}

// EventParticipantAdded описывает событие добавления участника
type EventParticipantAdded struct {
	events.Head
	ChatID      uuid.UUID
	Participant Participant
}

// EventParticipantRemoved описывает событие удаления участника
type EventParticipantRemoved struct {
	events.Head
	ChatID      uuid.UUID
	Participant Participant
}

// EventChatNameUpdated описывает событие обновления названия чата
type EventChatNameUpdated struct {
	events.Head
	ChatID uuid.UUID
	Name   string
}

// EventChatCreated описывает событие создания чата
type EventChatCreated struct {
	events.Head
	ChatID uuid.UUID
}
