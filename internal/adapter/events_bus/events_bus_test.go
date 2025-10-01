package eventsBus

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nice-pea/npchat/internal/usecases/events"
)

func Test_EventsBus(t *testing.T) {
	// id сессии + id пользователя уникальный ключ для
	t.Run("сессия не может начать второе прослушивание, если уже есть активное", func(t *testing.T) {
		b := new(EventsBus)

		sessionID := uuid.New()
		userID := uuid.New()

		_, err := b.AddListener(userID, sessionID, nil)
		require.NoError(t, err)

		_, err = b.AddListener(userID, sessionID, nil)
		require.ErrorIs(t, err, ErrDuplicateSession)
	})

	t.Run("после отмены прослушивания, можно подписаться заново", func(t *testing.T) {
		b := new(EventsBus)

		sessionID := uuid.New()
		userID := uuid.New()

		// Запустить прослушивание
		removeListener, err := b.AddListener(userID, sessionID, nil)
		require.NoError(t, err)

		// Отменить прослушивание со стороны подписчикаы
		removeListener()

		// Запустить прослушивание
		_, err = b.AddListener(userID, sessionID, func(event events.Event, err error) {
			// Ошибка об отмене прослушивания сервером
			require.ErrorIs(t, err, ErrListenerForciblyCanceled)
			require.Zero(t, event)
		})
		require.NoError(t, err)

		// Отменить прослушивание со стороны сервера
		b.Cancel(sessionID)

		// Запустить прослушивание
		_, err = b.AddListener(userID, sessionID, nil)
		require.NoError(t, err)
	})

	t.Run("healthcheck позволяет вытеснить неактивного слушателя", func(t *testing.T) {
		b := new(EventsBus)

		sessionID := uuid.New()
		userID := uuid.New()

		// Создать заблокированный слушатель (не обрабатывает события)
		_, err := b.AddListener(userID, sessionID, func(event events.Event, err error) {
			// Заблокировать обработчик навсегда
			select {}
		})
		require.NoError(t, err)

		// Попытаться добавить нового слушателя для той же сессии
		// Healthcheck должен обнаружить, что старый слушатель не отвечает
		eventReceived := make(chan bool, 1)
		_, err = b.AddListener(userID, sessionID, func(event events.Event, err error) {
			if event.Type == "test" {
				eventReceived <- true
			}
		})
		require.NoError(t, err)

		// Отправить тестовое событие
		b.Consume([]events.Event{{Type: "test", Recipients: []uuid.UUID{userID}}})

		// Убедиться, что новый слушатель получил событие
		select {
		case <-eventReceived:
			// Успех
		case <-time.After(time.Second):
			t.Fatal("новый слушатель не получил событие")
		}
	})

	t.Run("healthcheck не вытесняет активного слушателя", func(t *testing.T) {
		b := new(EventsBus)

		sessionID := uuid.New()
		userID := uuid.New()

		// Создать активный слушатель (быстро обрабатывает события)
		_, err := b.AddListener(userID, sessionID, func(event events.Event, err error) {
			// Быстро вернуться
		})
		require.NoError(t, err)

		// Попытаться добавить нового слушателя для той же сессии
		// Healthcheck должен определить, что старый слушатель активен
		_, err = b.AddListener(userID, sessionID, func(event events.Event, err error) {})
		require.ErrorIs(t, err, ErrDuplicateSession)
	})

	t.Run("в сессии будут приходить только события направляемые пользователям", func(t *testing.T) {
		b := new(EventsBus)
		// Счетчик событий
		eventsCountByUserID := map[uuid.UUID]int{}
		var mu sync.RWMutex
		// Пользователи
		userIDs := []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		}
		// Отправляемые события
		eventsForSend := []events.Event{
			events.Event{Recipients: []uuid.UUID{userIDs[0], userIDs[1], userIDs[2]}},
			events.Event{Recipients: []uuid.UUID{userIDs[1], userIDs[2]}},
			events.Event{Recipients: []uuid.UUID{userIDs[2]}},
		}

		// Запустить прослушивание
		for _, userID := range userIDs {
			_, err := b.AddListener(userID, userID, func(event events.Event, _ error) {
				mu.Lock()
				eventsCountByUserID[userID]++
				mu.Unlock()
			})
			require.NoError(t, err)
		}

		// Отправить события
		b.Consume(eventsForSend)

		// Проверить количество полученных событий каждым пользователем
		assert.Equal(t, 1, eventsCountByUserID[userIDs[0]])
		assert.Equal(t, 2, eventsCountByUserID[userIDs[1]])
		assert.Equal(t, 3, eventsCountByUserID[userIDs[2]])
	})

	t.Run("прослушивание можно отменять стороной шины", func(t *testing.T) {
		b := new(EventsBus)

		sessionID := uuid.New()
		userID := uuid.New()

		// Запустить прослушивание
		_, err := b.AddListener(userID, sessionID, func(event events.Event, err error) {})
		require.NoError(t, err)
		time.Sleep(time.Millisecond)

		// Отменить сессию
		b.Cancel(sessionID)
		// Проверить список просшушиваний
		assert.Empty(t, b.activeListeners())
	})

	t.Run("закрытие сервера удалит всех слушателей и не будет принимать новых", func(t *testing.T) {
		b := new(EventsBus)

		// Запустить много слушаетелй
		for range 100 {
			_, err := b.AddListener(uuid.New(), uuid.New(), func(event events.Event, err error) {})
			require.NoError(t, err)
		}

		// Закрыть сервер
		b.Close()
		// Попытаться запустить нового слушателя
		_, err := b.AddListener(uuid.New(), uuid.New(), nil)
		require.ErrorIs(t, err, ErrBusClosed)
		// Убедиться что список слушателей пуст
		assert.Empty(t, b.activeListeners())
	})

	t.Run("слушатель может отменить прослушивание", func(t *testing.T) {
		b := new(EventsBus)

		var removeListeners []func()

		receivedEvents := new(atomic.Int32)
		listenerIDs := lo.RepeatBy(100, func(_ int) uuid.UUID {
			return uuid.New()
		})

		// Запустить много слушаетелй
		for _, id := range listenerIDs {
			rl, err := b.AddListener(id, id, func(event events.Event, err error) {
				receivedEvents.Add(1)
			})
			require.NoError(t, err)
			removeListeners = append(removeListeners, rl)
		}

		// Отправить событие
		b.Consume([]events.Event{events.Event{Recipients: listenerIDs}})
		// Убедиться что все слушатели получили событие
		assert.Equal(t, len(listenerIDs), int(receivedEvents.Load()))

		// Удалить слушателей
		for _, removeListener := range removeListeners {
			removeListener()
		}

		// Проверить что список слsушателей пуст
		assert.Empty(t, b.activeListeners())

		// Отправить событие
		b.Consume([]events.Event{events.Event{Recipients: listenerIDs}})
		// Убедиться что никто не обработал событие
		assert.Equal(t, len(listenerIDs), int(receivedEvents.Load()))
	})
}
