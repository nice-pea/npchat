package chatt

import (
	"time"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

const (
	EventInvitationRemoved  = "invitation_removed"
	EventInvitationAdded    = "invitation_added"
	EventParticipantAdded   = "participant_added"
	EventParticipantRemoved = "participant_removed"
	EventChatCreated        = "chat_created"
	EventChatUpdated        = "chat_updated"
)

// NewEventInvitationRemoved описывает событие удаления приглашения
func (c *Chat) NewEventInvitationRemoved(invitation Invitation) events.Event {
	return events.Event{
		Type:      EventInvitationRemoved,
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
		Type:      EventInvitationAdded,
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
		Type:       EventParticipantAdded,
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
		Type:       EventParticipantRemoved,
		CreatedIn:  time.Now(),
		Recipients: append(userIDs(c.Participants), participant.UserID),
		Data: map[string]any{
			"chat":        *c,
			"participant": participant,
		},
	}
}

// NewEventChatCreated описывает событие создания чата
func (c *Chat) NewEventChatCreated() events.Event {
	return events.Event{
		Type:       EventChatCreated,
		CreatedIn:  time.Now(),
		Recipients: userIDs(c.Participants),
		Data: map[string]any{
			"chat": *c,
		},
	}
}

// NewEventChatUpdated описывает событие обновления чата
func (c *Chat) NewEventChatUpdated() events.Event {
	return events.Event{
		Type:       EventChatUpdated,
		CreatedIn:  time.Now(),
		Recipients: userIDs(c.Participants),
		Data: map[string]any{
			"chat": *c,
		},
	}
}
