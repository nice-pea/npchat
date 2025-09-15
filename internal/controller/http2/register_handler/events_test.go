package register_handler

import (
	"bufio"
	"context"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockRegisterHandler "github.com/nice-pea/npchat/internal/controller/http2/register_handler/mocks"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/usecases/events"
	findSession "github.com/nice-pea/npchat/internal/usecases/sessions/find_session"
)

var (
	mockSession = sessionn.Session{
		ID:     uuid.New(),
		UserID: uuid.New(),
	}
)

func TestEvents(t *testing.T) {
	// Создание моков
	mockSessionFinder := mockRegisterHandler.NewUsecasesForEvents(t)
	mockSessionFinder.
		On("FindSessions", mock.Anything).
		Return(findSession.Out{
			Sessions: []sessionn.Session{mockSession},
		}, nil)

	fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})

	// Регистрация обработчика
	Events(fiberApp, mockSessionFinder, mockEventListener{}, nil)

	// Запуск сервера на свободном локальном порту
	listener, err := net.Listen("tcp", ":")
	require.NoError(t, err)
	go func() { assert.NoError(t, fiberApp.Listener(listener)) }()
	// Отложить остановку сервера
	defer assert.NoError(t, fiberApp.Shutdown())

	// Контекст для корректной отмены запроса
	ctx, cancel := context.WithCancel(context.Background())

	// Создать запрос
	req, err := http.NewRequestWithContext(ctx, "GET", "http://"+listener.Addr().String()+"/events", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "SessionToken 123")

	// Выполнить запрос
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer func() { _ = resp.Body.Close() }()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Отменить запрос через некоторое время
	time.AfterFunc(time.Millisecond*5, cancel)

	// Читаем данные из потока.
	// Каждая итерация выполняется после вызова w.Flush() в обработчике
	var receivedEvents int
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		receivedEvents++
	}
	assert.Greater(t, receivedEvents, 0)
	log.Printf("получено событий: %d", receivedEvents)

	// На этом моменте запрос остановлен контекстом
	assert.False(t, scanner.Scan())
	assert.ErrorIs(t, scanner.Err(), context.Canceled)
}

type mockEventListener struct{}

func (m mockEventListener) AddListener(userID, sessionID uuid.UUID, f func(event events.Event, err error)) (removeListener func(), err error) {
	go func() {
		for {
			f(rndDirtyEvent(), nil)
			time.Sleep(time.Millisecond)
		}
	}()
	return func() {}, nil
}

func rndDirtyEvent() events.Event {
	return events.Event{
		Type:       "dirty",
		CreatedIn:  time.Now(),
		Recipients: []uuid.UUID{uuid.New()},
		Data: map[string]any{
			"name": gofakeit.Name(),
			"age":  gofakeit.Number(1, 100),
			"sex":  gofakeit.Gender(),
		},
	}
}
