package events

import (
	"time"

	"github.com/google/uuid"
)

type Head interface {
	CreatedIn() time.Time
	Recipients() []uuid.UUID
}

func NewHead(recipients []uuid.UUID) Head {
	return head{
		createdIn:  time.Now(),
		recipients: recipients,
	}
}

// head описывает событие
type head struct {
	createdIn  time.Time
	recipients []uuid.UUID
}

func (h head) CreatedIn() time.Time {
	return h.createdIn
}

func (h head) Recipients() []uuid.UUID {
	return h.recipients
}
