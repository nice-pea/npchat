package eventsBus

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EventsBus(t *testing.T) {
	// id сессии + id пользователя уникальный ключ для
	t.Run("сессия не может начать второе прослушивание, если уже есть активное", func(t *testing.T) {
		b := new(EventsBus)

		sessionID := uuid.New()
		userID := uuid.New()

		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			err := b.Listen(ctx, userID, sessionID, func(event any) {})
			require.ErrorIs(t, err, context.Canceled)
		}()
		time.Sleep(time.Millisecond)

		err := b.Listen(context.Background(), userID, sessionID, func(event any) {})
		require.Error(t, err)
		cancel()
	})

	t.Run("в сессии будут приходить только события направляемые пользователям", func(t *testing.T) {
		b := new(EventsBus)
		// Счетчик событий
		eventsCountByUserID := map[uuid.UUID]int{}
		// Пользователи
		userIDs := []uuid.UUID{
			uuid.New(),
			uuid.New(),
			uuid.New(),
		}
		// Отправляемые события
		events := []any{
			dummyEvent{recipients: []uuid.UUID{userIDs[0], userIDs[1], userIDs[2]}},
			dummyEvent{recipients: []uuid.UUID{userIDs[1], userIDs[2]}},
			dummyEvent{recipients: []uuid.UUID{userIDs[2]}},
		}

		// Запустить прослушивание
		ctx, cancel := context.WithCancel(context.Background())
		for _, userID := range userIDs {
			go b.Listen(ctx, userID, userID, func(event any) {
				eventsCountByUserID[userID]++
			})
		}
		time.Sleep(time.Millisecond)

		// Отправить события
		b.Consume(events)
		// Отменить прослушивание
		cancel()

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
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := b.Listen(context.Background(), userID, sessionID, func(event any) {})
			require.Error(t, err)
		}()
		time.Sleep(time.Millisecond)

		// Отменить сессию
		b.Cancel(sessionID)
		// Проверить список просшушиваний
		assert.Empty(t, b.listeners)
		wg.Wait()
	})

	t.Run("закрытие сервера отменить все прослушивания и не будет принимать новые", func(t *testing.T) {
		b := new(EventsBus)

		// Запустить много слушаетелй
		wg := sync.WaitGroup{}
		wg.Add(100)
		for range 100 {
			go func() {
				defer wg.Done()
				err := b.Listen(context.Background(), uuid.New(), uuid.New(), func(event any) {})
				require.Error(t, err)
			}()
		}

		// Закрыть сервер
		b.Close()
		// Попытаться запустить нового слушателя
		err := b.Listen(context.Background(), uuid.New(), uuid.New(), func(event any) {})
		require.Error(t, err)
		// Дождаться когда все слушатели завершат работу
		wg.Wait()
		// Убедиться что список слушателей пуст
		assert.Empty(t, b.listeners)
	})
}

type dummyEvent struct {
	recipients []uuid.UUID
}

func (e dummyEvent) CreatedIn() time.Time {
	return time.Time{}
}

func (e dummyEvent) Recipients() []uuid.UUID {
	return e.recipients
}
