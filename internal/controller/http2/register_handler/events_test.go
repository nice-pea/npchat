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
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	mockRegisterHandler "github.com/nice-pea/npchat/internal/controller/http2/register_handler/mocks"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	findSession "github.com/nice-pea/npchat/internal/usecases/sessions/find_session"
)

var (
	mockSession = sessionn.Session{
		ID:     uuid.New(),
		Status: sessionn.StatusNew,
		UserID: uuid.New(),
		Name:   "name",
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
	Events(fiberApp, mockSessionFinder, mockEventListener{})

	// Запуск сервера
	go func() { assert.NoError(t, fiberApp.Listen("localhost:8418")) }()
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
	assert.Greater(t, receivedEvents, 0)
	log.Printf("получено событий: %d", receivedEvents)

	// На этом моменте запрос остановлен контекстом
	assert.False(t, scanner.Scan())
	assert.ErrorIs(t, scanner.Err(), context.Canceled)
	// Остановить сервер
	assert.NoError(t, fiberApp.Shutdown())

}

type mockEventListener struct{}

func (m mockEventListener) Listen(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, eventHandler func(event any)) error {
	for {
		select {
		case <-ctx.Done():
			log.Print("Event Listener получил сигнал завершения контекста из http-обработчика")
			return ctx.Err()
		default:
			eventHandler(rndDirtyEvent())
		}
	}
}

func rndDirtyEvent() any {
	var e struct {
		Name string `fake:"{name}"`
		Age  int    `fake:"{number:0,100}"`
		Sex  string `fake:"{gender}"`
	}
	gofakeit.Struct(&e)
	return e
}
