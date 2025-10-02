package eventsBus

import (
	"errors"
	"slices"
	"sync"
	"sync/atomic"

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
	userID    uuid.UUID                           // ID пользователя
	sessionID uuid.UUID                           // ID сессии
	f         func(event events.Event, err error) // Обработчик событий
}

// AddListener регистрирует обработчик событий
// При конфликте (дублирующая сессия) применяется принцип FIFO: первый подписчик принудительно отписывается
func (u *EventsBus) AddListener(userID, sessionID uuid.UUID, f func(event events.Event, err error)) (removeListener func(), err error) {
	// Проверить, что сервер не закрыт
	if u.closed.Load() {
		return nil, ErrBusClosed
	}

	u.listenersMutex.Lock()

	// Найти существующего слушателя для этой сессии
	existingListenerIndex := slices.IndexFunc(u.listeners, func(l *listener) bool {
		return l.userID == userID && l.sessionID == sessionID
	})

	// Если найден существующий слушатель, применяем принцип FIFO
	if existingListenerIndex != -1 {
		existingListener := u.listeners[existingListenerIndex]

		// Удалить старого слушателя из списка
		u.listeners = slices.Delete(u.listeners, existingListenerIndex, existingListenerIndex+1)

		// Отправить ошибку старому слушателю в отдельной горутине
		if existingListener.f != nil {
			go func() {
				defer func() {
					// Восстановиться от любой паники
					_ = recover()
				}()
				existingListener.f(events.Event{}, ErrListenerForciblyCanceled)
			}()
		}
	}

	listener := &listener{
		userID:    userID,
		sessionID: sessionID,
		f:         f,
	}
	// Добавить слушателя
	u.listeners = append(u.listeners, listener)

	u.listenersMutex.Unlock()

	return func() {
		u.listenersMutex.Lock()
		defer u.listenersMutex.Unlock()
		// Удалить слушателя из списка
		u.listeners = slices.DeleteFunc(u.listeners, func(l *listener) bool {
			return l == listener
		})
	}, nil
}


// Consume рассылает события слушателям
func (u *EventsBus) Consume(ee []events.Event) {
	// Выйти, если сервер уже закрыт
	if u.closed.Load() {
		return
	}

	// Снимок слушателей
	u.listenersMutex.Lock()
	snapshot := make([]*listener, len(u.listeners))
	copy(snapshot, u.listeners)
	u.listenersMutex.Unlock()

	// Выйти, если нет слушателей
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
			if packet.listener.f != nil {
				packet.listener.f(packet.event, nil)
			}
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

	// Снимок слушателей
	u.listenersMutex.Lock()
	snapshot := make([]*listener, len(u.listeners))
	copy(snapshot, u.listeners)
	u.listenersMutex.Unlock()

	// Инициализация waitgroup для асинхронной отмены
	var wg sync.WaitGroup
	wg.Add(len(snapshot))

	// Отправить ошибку всем активным слушателям
	for _, listener := range snapshot {
		go func() {
			defer wg.Done()
			if listener.f != nil {
				listener.f(events.Event{}, ErrBusClosed)
			}
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
		return l.sessionID == sessionID
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
	if target.f != nil {
		target.f(events.Event{}, ErrListenerForciblyCanceled)
	}
}

