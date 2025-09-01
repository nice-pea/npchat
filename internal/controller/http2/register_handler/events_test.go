package register_handler

import (
	"bufio"
	"context"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nice-pea/npchat/internal/domain/sessionn"
	findSession "github.com/nice-pea/npchat/internal/usecases/sessions/find_session"
)

func TestEvents(t *testing.T) {
	server := fiber.New(fiber.Config{
		ReadTimeout:           time.Second,
		WriteTimeout:          time.Second,
		DisableKeepalive:      false,
		DisableStartupMessage: true,
		StreamRequestBody:     true,
	})
	// Регистрация обработчика
	Events(server, mockSessionFinder{}, mockEventListener{})

	// Запуск сервера
	go func() { assert.NoError(t, server.Listen("localhost:8418")) }()
	// Задержка чтобы сервер успел запуститься
	time.Sleep(time.Millisecond * 5)

	// Контекст для корректной отмены запроса
	ctx, cancel := context.WithCancel(context.Background())

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8418/events", nil)
	require.NoError(t, err)

	// Выполнить запрос
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Отменить запрос через некоторое время
	time.AfterFunc(time.Millisecond*5, func() {
		cancel()
	})

	var receivedEvents int

	// Читаем данные из потока.
	// Каждая итерация выполняется после вызова w.Flush() в обработчике
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		receivedEvents++
	}

	// На этом моменте запрос остановлен контекстом
	assert.False(t, scanner.Scan())
	assert.ErrorIs(t, scanner.Err(), context.Canceled)
	// Остановить сервер
	assert.NoError(t, server.Shutdown())
	log.Printf("получено событий: %d", receivedEvents)
}

var (
	mockSession = sessionn.Session{
		ID:     uuid.New(),
		Status: sessionn.StatusNew,
		UserID: uuid.New(),
		Name:   "name",
	}
)

type mockSessionFinder struct{}

func (m mockSessionFinder) FindSessions(findSession.In) (findSession.Out, error) {
	return findSession.Out{
		Sessions: []sessionn.Session{mockSession},
	}, nil
}

type mockEventListener struct{}

func (m mockEventListener) Listen(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, eventHandler func(event any)) error {
	type someEventData struct {
		Name string `fake:"{name}"`
		Age  int    `fake:"{number:0,100}"`
		Sex  string `fake:"{gender}"`
	}

	for {
		select {
		case <-ctx.Done():
			log.Print("Event Listener получил сигнал завершения контекста из http-обработчика")
			return ctx.Err()
		default:
			var e someEventData
			gofakeit.Struct(&e)
			eventHandler(e)
		}
	}
}
