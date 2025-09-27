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

const (
	healthcheckTimeout = 100 * time.Millisecond
)

var (
	// ErrBusClosed возвращается при попытке добавить слушателя к закрытой шине событий
	ErrBusClosed = errors.New("шина событий закрыта")
	// ErrDuplicateSession возвращается при попытке создать дублирующуюся подписку для активной сессии
	ErrDuplicateSession = errors.New("сессия уже прослушивает события")
	// ErrListenerForciblyCanceled отправляется слушателю при принудительной отмене подписки
	ErrListenerForciblyCanceled = errors.New("принудительно отменен")
)

// EventsBus представляет шину событий для управления подписчиками и рассылки событий.
// Обеспечивает потокобезопасное добавление, удаление слушателей и отправку событий.
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

// AddListener регистрирует обработчик событий для указанного пользователя и сессии.
// При обнаружении существующей подписки с той же сессией выполняется healthcheck:
// - если существующий слушатель не отвечает (неактивен), он удаляется и регистрируется новый
// - если существующий слушатель активен, возвращается ErrDuplicateSession
// Возвращает функцию для отмены подписки и возможную ошибку.
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
		if existingListener.f == nil {
			return nil, ErrDuplicateSession
		}
		// Отпустить мьютекс для healthcheck, чтобы избежать deadlock
		u.listenersMutex.Unlock()
		healthy := u.healthcheck(existingListener)
		u.listenersMutex.Lock()

		// Повторно проверить, что слушатель всё ещё существует
		stillExists := slices.IndexFunc(u.listeners, func(l *listener) bool {
			return l.sessionID == sessionID && !l.listeningIsOver
		}) != -1

		if healthy || !stillExists {
			// Слушатель активен или уже удален, вернуть ошибку дубликата
			return nil, ErrDuplicateSession
		}

		// Слушатель не активен, собрать всех удаляемых для уведомления
		toCancel := lo.Filter(u.listeners, func(l *listener, _ int) bool {
			return l.sessionID == sessionID && !l.listeningIsOver
		})

		// Удалить всех по сессии
		u.listeners = slices.DeleteFunc(u.listeners, func(l *listener) bool {
			return l.sessionID == sessionID
		})

		// Отправить ошибку всем удаленным слушателям в отдельных горутинах (безопасно)
		for _, l := range toCancel {
			if l.f != nil {
				go func(fn func(event events.Event, err error)) {
					defer func() { recover() }()
					fn(events.Event{}, ErrListenerForciblyCanceled)
				}(l.f)
			}
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

// Consume рассылает массив событий соответствующим слушателям.
// Каждое событие отправляется только тем слушателям, чьи userID указаны в списке получателей.
// Рассылка выполняется асинхронно с ожиданием завершения всех отправок.
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
		p := packet
		go func(pkt forHandling) {
			defer wg.Done()
			defer func() { recover() }()
			if pkt.listener.f != nil {
				pkt.listener.f(pkt.event, nil)
			}
		}(p)
	}

	// Ожидать завершения обработки событий слушателями
	wg.Wait()
}

// Close завершает работу шины событий, отправляя ErrBusClosed всем активным слушателям
// и очищая список подписчиков. После вызова Close новые подписки не принимаются.
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
		l := listener
		go func(lst *listener) {
			defer wg.Done()
			defer func() { recover() }()
			if lst.f != nil {
				lst.f(events.Event{}, ErrBusClosed)
			}
		}(l)
	}

	// Ожидать завершения обработки ошибки слушателями
	wg.Wait()

	// Очистить список
	u.listenersMutex.Lock()
	u.listeners = nil
	u.listenersMutex.Unlock()
}

// Cancel принудительно отменяет подписку слушателя по sessionID,
// отправляя ему ErrListenerForciblyCanceled и удаляя из списка активных подписчиков.
func (u *EventsBus) Cancel(sessionID uuid.UUID) {
	// Выйти, если сервер уже закрыт
	if u.closed.Load() {
		return
	}

	u.listenersMutex.Lock()

	// Собрать всех удаляемых слушателей для уведомления
	toCancel := lo.Filter(u.listeners, func(l *listener, _ int) bool {
		return l.sessionID == sessionID && !l.listeningIsOver
	})

	if len(toCancel) == 0 {
		u.listenersMutex.Unlock()
		return
	}

	// Удалить всех слушателей по сессии
	u.listeners = slices.DeleteFunc(u.listeners, func(l *listener) bool {
		return l.sessionID == sessionID
	})

	u.listenersMutex.Unlock()

	// Отправить ошибку всем удаленным слушателям (безопасно)
	for _, l := range toCancel {
		if l.f != nil {
			go func(fn func(event events.Event, err error)) {
				defer func() { recover() }()
				fn(events.Event{}, ErrListenerForciblyCanceled)
			}(l.f)
		}
	}
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
	if l == nil || l.f == nil {
		return true
	}
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
	timer := time.NewTimer(healthcheckTimeout)
	defer timer.Stop()

	select {
	case <-ack:
		// Слушатель ответил, он активен
		return true
	case <-timer.C:
		// Таймаут, слушатель не активен
		return false
	}
}
