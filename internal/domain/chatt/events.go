package chatt

import (
	"time"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

// EventInvitationRemoved описывает событие удаления приглашения

const (
	EventInvitationRemovedType  = "invitation_removed"
	EventInvitationAddedType    = "invitation_added"
	EventParticipantAddedType   = "participant_added"
	EventParticipantRemovedType = "participant_removed"
	EventChatNameUpdatedType    = "chat_name_updated"
	EventChatCreatedType        = "chat_created"
)

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
			"invitation": invitation,
		},
	}
}

// EventInvitationAdded описывает событие добавления приглашения
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
			"invitation": invitation,
		},
	}
}

// EventParticipantAdded описывает событие добавления участника
func (c *Chat) NewEventParticipantAdded(participant Participant) events.Event {
	return events.Event{
		Type:       EventParticipantAddedType,
		CreatedIn:  time.Now(),
		Recipients: userIDs(c.Participants),
		Data: map[string]any{
			"chat_id":     c.ID,
			"participant": participant,
		},
	}
}

// EventParticipantRemoved описывает событие удаления участника
func (c *Chat) NewEventParticipantRemoved(participant Participant) events.Event {
	return events.Event{
		Type:       EventParticipantRemovedType,
		CreatedIn:  time.Now(),
		Recipients: append(userIDs(c.Participants), participant.UserID),
		Data: map[string]any{
			"chat_id":     c.ID,
			"participant": participant,
		},
	}
}

// EventChatNameUpdated описывает событие обновления названия чата
func (c *Chat) NewEventChatNameUpdated() events.Event {
	return events.Event{
		Type:       EventChatNameUpdatedType,
		CreatedIn:  time.Now(),
		Recipients: userIDs(c.Participants),
		Data: map[string]any{
			"chat_id": c.ID,
			"name":    c.Name,
		},
	}
}

// EventChatCreated описывает событие создания чата
func (c *Chat) NewEventChatCreated() events.Event {
	return events.Event{
		Type:       EventChatCreatedType,
		CreatedIn:  time.Now(),
		Recipients: userIDs(c.Participants),
		Data: map[string]any{
			"chat_id": c.ID,
		},
	}
}
