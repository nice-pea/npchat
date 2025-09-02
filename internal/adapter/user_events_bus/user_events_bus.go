package userEventsBus

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

type UserEventsBus struct {
	listeners []listener
}

type listener struct {
	err       chan<- error
	userID    uuid.UUID
	sessionID uuid.UUID
	f         func(event any)
}

func (u *UserEventsBus) Listen(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, f func(event any)) error {
	// Проверить, что слушатель ещё не зарегистрирован
	sessionAlreadyListen := slices.ContainsFunc(u.listeners, func(l listener) bool {
		return l.userID == userID && l.sessionID == sessionID
	})
	if sessionAlreadyListen {
		return errors.New("session already listen events")
	}

	errChan := make(chan error)
	defer close(errChan)

	// Добавить слушателя
	u.listeners = append(u.listeners, listener{
		// ctx:       ctx,
		userID:    userID,
		sessionID: sessionID,
		f:         f,
		err:       errChan,
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func (u *UserEventsBus) Consume(ee []any) {
	for _, event := range ee {
		eventHead, ok := event.(events.Head)
		if !ok {
			continue
		}

		// Найти получателей события
		recipients := slices.DeleteFunc(u.listeners, func(l listener) bool {
			return !slices.ContainsFunc(eventHead.Recipients(), func(userID uuid.UUID) bool {
				return l.userID == userID
			})
		})

		// Отправить событие получателям
		for _, r := range recipients {
			r.f(event)
		}
	}
}

func (u *UserEventsBus) Close() {
	for _, l := range u.listeners {
		l.err <- errors.New("сервер закрыт")
	}
}

func (u *UserEventsBus) Cancel(sessionID uuid.UUID) {
	i := slices.IndexFunc(u.listeners, func(l listener) bool {
		return l.sessionID == sessionID
	})
	if i == -1 {
		return
	}

	u.listeners[i].err <- errors.New("принудительно отменен")

	u.listeners = slices.DeleteFunc(u.listeners, func(l listener) bool {
		return l.sessionID == sessionID
	})
}
