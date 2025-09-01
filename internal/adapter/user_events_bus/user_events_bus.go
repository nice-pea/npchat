package userEventsBus

import (
	"context"

	"github.com/google/uuid"
)

type UserEventsBus struct{}

func (u *UserEventsBus) Listen(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, f func(event any)) error {
	panic("implement me")
}
