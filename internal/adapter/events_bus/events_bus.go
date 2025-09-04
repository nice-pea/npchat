package eventsBus

import (
	"errors"
	"slices"
	"sync"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

type EventsBus struct {
	listeners []*listener
	closed    bool
	mu        sync.RWMutex
}

type listener struct {
	listeningIsOver bool
	userID          uuid.UUID
	sessionID       uuid.UUID
	f               func(event any, err error)
}

func (u *EventsBus) AddListener(userID, sessionID uuid.UUID, f func(event any, err error)) (removeListener func(), err error) {
	// Проверить, что сервер не закрыт
	if u.closed {
		return nil, errors.New("events bus is closed")
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	// Проверить, что слушатель ещё не зарегистрирован
	sessionAlreadyListen := slices.ContainsFunc(u.listeners, func(l *listener) bool {
		return l.userID == userID && l.sessionID == sessionID
	})
	if sessionAlreadyListen {
		return nil, errors.New("session already listen events")
	}

	listener := &listener{
		userID:    userID,
		sessionID: sessionID,
		f:         f,
	}
	// Добавить слушателя
	u.listeners = append(u.listeners, listener)

	return func() {
		// Пометить слушателя как окончившего слушать события
		listener.listeningIsOver = true
	}, nil
}

func (u *EventsBus) Consume(ee []any) {
	// Выйти, если сервер уже закрыт
	if u.closed {
		return
	}

	u.mu.Lock()
	defer u.mu.Unlock()

	// Очистить слушателей, которые отменили подписки
	u.listeners = slices.DeleteFunc(u.listeners, func(l *listener) bool {
		return l.listeningIsOver
	})

	// Пара событие + слушатель, для удобной отправки
	type forHandling struct {
		listener *listener
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
		recipients := lo.Filter(u.listeners, func(l *listener, _ int) bool {
			return slices.ContainsFunc(eventHead.Recipients(), func(userID uuid.UUID) bool {
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
			packet.listener.f(packet.event, nil)
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

	// Инициализация waitgroup для асинхронной отмены
	var wg sync.WaitGroup
	wg.Add(len(u.listeners))

	// Отправить ошибку всем слушателям
	for _, listener := range u.listeners {
		go func() {
			defer wg.Done()
			listener.f(nil, errors.New("сервер закрыт"))
		}()
	}

	// Ожидать завершения обработки ошибки слушателями
	wg.Wait()

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
	i := slices.IndexFunc(u.listeners, func(l *listener) bool {
		return l.sessionID == sessionID && !l.listeningIsOver
	})
	if i == -1 {
		return
	}

	// Отправить ошибку
	u.listeners[i].f(nil, errors.New("принудительно отменен"))

	// Удалить слушателя
	u.listeners = slices.DeleteFunc(u.listeners, func(l *listener) bool {
		return l.sessionID == sessionID
	})
}

func (u *EventsBus) activeListeners() []*listener {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return slices.DeleteFunc(u.listeners, func(l *listener) bool {
		return l.listeningIsOver
	})
}
