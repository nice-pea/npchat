package eventsBus

import (
	"context"
	"errors"
	"slices"
	"sync"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

type EventsBus struct {
	listeners []listener
	closed    bool
	mu        sync.RWMutex
}

type listener struct {
	err       chan<- error
	userID    uuid.UUID
	sessionID uuid.UUID
	f         func(event any)
}

func (u *EventsBus) Listen(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, f func(event any)) error {
	// Проверить, что сервер не закрыт
	if u.closed {
		return errors.New("events bus is closed")
	}
	u.mu.Lock()
	// Проверить, что слушатель ещё не зарегистрирован
	sessionAlreadyListen := slices.ContainsFunc(u.listeners, func(l listener) bool {
		return l.userID == userID && l.sessionID == sessionID
	})
	if sessionAlreadyListen {
		u.mu.Unlock()
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
	u.mu.Unlock()

	select {
	case <-ctx.Done():
		u.mu.Lock()
		// Удалить слушателя
		u.listeners = slices.DeleteFunc(u.listeners, func(l listener) bool {
			return l.sessionID == sessionID
		})
		u.mu.Unlock()

		return ctx.Err()
	case err := <-errChan:
		return err
	}
}

func (u *EventsBus) Consume(ee []any) {
	// Выйти, если сервер уже закрыт
	if u.closed {
		return
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	// Пара событие + слушатель, для удобной отправки
	type forHandling struct {
		listener listener
		event    any
	}

	// Одномерный список из события и их получателей
	var le []forHandling

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

		// Собрать события и их получателей к отправке
		for _, r := range recipients {
			le = append(le, forHandling{listener: r, event: event})
		}
	}

	// Инициализация waitgroup для асинхронной отправки
	var wg sync.WaitGroup
	wg.Add(len(le))

	// Отправить событие получателям
	for _, packet := range le {
		go func() {
			defer wg.Done()
			packet.listener.f(packet.event)
		}()
	}

	// Ожидать завершения обработки событий слушателями
	wg.Wait()
}

func (u *EventsBus) Close() {
	// Выйти, если сервер уже закрыт
	if u.closed {
		return
	}

	// Установить флаг closed
	u.closed = true

	u.mu.Lock()
	defer u.mu.Unlock()

	// Отправить ошибку всем слушателям
	for _, l := range u.listeners {
		l.err <- errors.New("сервер закрыт")
	}
	// Очистить список
	u.listeners = nil
}

func (u *EventsBus) Cancel(sessionID uuid.UUID) {
	// Выйти, если сервер уже закрыт
	if u.closed {
		return
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	// Найти слушателя по сессии
	i := slices.IndexFunc(u.listeners, func(l listener) bool {
		return l.sessionID == sessionID
	})
	if i == -1 {
		return
	}

	// Отправить ошибку
	u.listeners[i].err <- errors.New("принудительно отменен")

	// Удалить слушателя
	u.listeners = slices.DeleteFunc(u.listeners, func(l listener) bool {
		return l.sessionID == sessionID
	})
}
