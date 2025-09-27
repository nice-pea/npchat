package eventsBus

import (
	"errors"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

var (
	ErrBusClosed                = errors.New("шина событий закрыта")
	ErrDuplicateSession         = errors.New("сессия уже прослушивает события")
	ErrListenerForciblyCanceled = errors.New("принудительно отменен")
)

type EventsBus struct {
	listeners      []*listener // Список слушателей
	listenersMutex sync.Mutex  // Синхронизация доступа к listeners
	closed         atomic.Bool // Признак окончания работы системы
}

// listener представляет собой слушателя (подписчика) событий
type listener struct {
	listeningIsOver bool                                // Признак отмены подписки, со стороны слушателя
	userID          uuid.UUID                           // ID пользователя
	sessionID       uuid.UUID                           // ID сессии
	f               func(event events.Event, err error) // Обработчик событий
}

// AddListener регистрирует обработчик событий
func (u *EventsBus) AddListener(userID, sessionID uuid.UUID, f func(event events.Event, err error)) (removeListener func(), err error) {
	// Проверить, что сервер не закрыт
	if u.closed.Load() {
		return nil, ErrBusClosed
	}

	u.listenersMutex.Lock()
	defer u.listenersMutex.Unlock()

	// Проверить, что слушатель ещё не зарегистрирован
	existingListenerIndex := slices.IndexFunc(u.listeners, func(l *listener) bool {
		return l.userID == userID &&
			l.sessionID == sessionID &&
			!l.listeningIsOver
	})
	if existingListenerIndex != -1 {
		existingListener := u.listeners[existingListenerIndex]
		// Выполнить healthcheck существующего слушателя
		if !u.healthcheck(existingListener) {
			// Слушатель не активен, принудительно отменить его
			u.listeners = slices.DeleteFunc(u.listeners, func(l *listener) bool {
				return l.sessionID == sessionID
			})
			// Отправить ошибку старому слушателю в отдельной горутине
			go existingListener.f(events.Event{}, ErrListenerForciblyCanceled)
		} else {
			// Слушатель активен, вернуть ошибку дубликата
			return nil, ErrDuplicateSession
		}
	}

	listener := &listener{
		userID:    userID,
		sessionID: sessionID,
		f:         f,
	}
	// Добавить слушателя
	u.listeners = append(u.listeners, listener)

	return func() {
		u.listenersMutex.Lock()
		// Пометить слушателя как окончившего слушать события
		listener.listeningIsOver = true
		u.listenersMutex.Unlock()
	}, nil
}

// Consume рассылает события слушателям
func (u *EventsBus) Consume(ee []events.Event) {
	// Выйти, если сервер уже закрыт
	if u.closed.Load() {
		return
	}

	// Снимок активных слушателей
	snapshot := u.activeListeners()

	// Выйти, если нет активных слушателей
	if len(snapshot) == 0 {
		return
	}

	// Пара событие + слушатель, для удобной отправки
	type forHandling struct {
		listener *listener
		event    events.Event
	}

	// Одномерный список из события и их получателей
	var le []forHandling

	for _, event := range ee {
		// Найти получателей события
		recipients := lo.Filter(snapshot, func(l *listener, _ int) bool {
			return slices.ContainsFunc(event.Recipients, func(userID uuid.UUID) bool {
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

// Close завершает работу системы
func (u *EventsBus) Close() {
	// Выйти, если сервер уже закрыт
	if !u.closed.CompareAndSwap(false, true) {
		return
	}

	// Снимок активных слушателей
	snapshot := u.activeListeners()

	// Инициализация waitgroup для асинхронной отмены
	var wg sync.WaitGroup
	wg.Add(len(snapshot))

	// Отправить ошибку всем активным слушателям
	for _, listener := range snapshot {
		go func() {
			defer wg.Done()
			listener.f(events.Event{}, ErrBusClosed)
		}()
	}

	// Ожидать завершения обработки ошибки слушателями
	wg.Wait()

	// Очистить список
	u.listenersMutex.Lock()
	u.listeners = nil
	u.listenersMutex.Unlock()
}

// Cancel отменяет подписку слушателя (удаляет его из списка)
func (u *EventsBus) Cancel(sessionID uuid.UUID) {
	// Выйти, если сервер уже закрыт
	if u.closed.Load() {
		return
	}

	u.listenersMutex.Lock()

	// Найти слушателя по сессии
	i := slices.IndexFunc(u.listeners, func(l *listener) bool {
		return l.sessionID == sessionID && !l.listeningIsOver
	})
	if i == -1 {
		u.listenersMutex.Unlock()
		return
	}
	target := u.listeners[i]

	// Удалить слушателя
	u.listeners = slices.DeleteFunc(u.listeners, func(l *listener) bool {
		return l.sessionID == sessionID
	})

	u.listenersMutex.Unlock()

	// Отправить ошибку
	target.f(events.Event{}, ErrListenerForciblyCanceled)
}

// activeListeners возвращает активных слушателей
func (u *EventsBus) activeListeners() []*listener {
	u.listenersMutex.Lock()
	defer u.listenersMutex.Unlock()
	return lo.Filter(u.listeners, func(l *listener, _ int) bool {
		return !l.listeningIsOver
	})
}

// healthcheck проверяет, является ли слушатель активным
// отправляя ему тестовое событие с таймаутом
func (u *EventsBus) healthcheck(l *listener) bool {
	// Создать канал для подтверждения получения
	ack := make(chan struct{}, 1)

	// Попытаться отправить тестовое событие слушателю
	go func() {
		defer func() {
			// Восстановиться от паники, если слушатель уже закрыт
			recover()
		}()
		l.f(events.Event{}, nil)
		ack <- struct{}{}
	}()

	// Ждать подтверждения с таймаутом
	select {
	case <-ack:
		// Слушатель ответил, он активен
		return true
	case <-time.After(100 * time.Millisecond):
		// Таймаут, слушатель не активен
		return false
	}
}
