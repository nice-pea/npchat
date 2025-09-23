package chatt

import (
	"time"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

const (
	EventInvitationRemovedType  = "invitation_removed"
	EventInvitationAddedType    = "invitation_added"
	EventParticipantAddedType   = "participant_added"
	EventParticipantRemovedType = "participant_removed"
	EventChatNameUpdatedType    = "chat_name_updated"
	EventChatCreatedType        = "chat_created"
	EventChatActiveUpdated      = "chat_active_updated"
)

// NewEventInvitationRemoved описывает событие удаления приглашения
func (c *Chat) NewEventInvitationRemoved(invitation Invitation) events.Event {
	return events.Event{
		Type:      EventInvitationRemovedType,
		CreatedIn: time.Now(),
		Recipients: []uuid.UUID{
			c.ChiefID,
			invitation.SubjectID,
			invitation.RecipientID,
		},
		Data: map[string]any{
			"chat":       *c,
			"invitation": invitation,
		},
	}
}

// NewEventInvitationAdded описывает событие добавления приглашения
func (c *Chat) NewEventInvitationAdded(invitation Invitation) events.Event {
	return events.Event{
		Type:      EventInvitationAddedType,
		CreatedIn: time.Now(),
		Recipients: []uuid.UUID{
			c.ChiefID,
			invitation.SubjectID,
			invitation.RecipientID,
		},
		Data: map[string]any{
			"chat":       *c,
			"invitation": invitation,
		},
	}
}

// NewEventParticipantAdded описывает событие добавления участника
func (c *Chat) NewEventParticipantAdded(participant Participant) events.Event {
	return events.Event{
		Type:       EventParticipantAddedType,
		CreatedIn:  time.Now(),
		Recipients: userIDs(c.Participants),
		Data: map[string]any{
			"chat":        *c,
			"participant": participant,
		},
	}
}

// NewEventParticipantRemoved описывает событие удаления участника
func (c *Chat) NewEventParticipantRemoved(participant Participant) events.Event {
	return events.Event{
		Type:       EventParticipantRemovedType,
		CreatedIn:  time.Now(),
		Recipients: append(userIDs(c.Participants), participant.UserID),
		Data: map[string]any{
			"chat":        *c,
			"participant": participant,
		},
	}
}

// NewEventChatNameUpdated описывает событие обновления названия чата
func (c *Chat) NewEventChatNameUpdated() events.Event {
	return events.Event{
		Type:       EventChatNameUpdatedType,
		CreatedIn:  time.Now(),
		Recipients: userIDs(c.Participants),
		Data: map[string]any{
			"chat": *c,
		},
	}
}

// NewEventChatCreated описывает событие создания чата
func (c *Chat) NewEventChatCreated() events.Event {
	return events.Event{
		Type:       EventChatCreatedType,
		CreatedIn:  time.Now(),
		Recipients: userIDs(c.Participants),
		Data: map[string]any{
			"chat": *c,
		},
	}
}

// NewEventActiveUpdated описывает событие обновления активности чата
func (c *Chat) NewEventActiveUpdated() events.Event {
	return events.Event{
		Type:       EventChatActiveUpdated,
		CreatedIn:  time.Now(),
		Recipients: userIDs(c.Participants),
		Data: map[string]any{
			"chat": *c,
		},
	}
}
