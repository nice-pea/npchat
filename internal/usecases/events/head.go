package events

import (
	"time"

	"github.com/google/uuid"
)

type Head interface {
	CreatedIn() time.Time
	Recipients() []uuid.UUID
}

func NewHead(recipients []uuid.UUID) head {
	return head{
		createdIn:  time.Now(),
		cecipients: recipients,
	}
}

// head описывает событие
type head struct {
	createdIn  time.Time
	cecipients []uuid.UUID
}

func (h head) CreatedIn() time.Time {
	return h.createdIn
}

func (h head) Recipients() []uuid.UUID {
	return h.cecipients
}
