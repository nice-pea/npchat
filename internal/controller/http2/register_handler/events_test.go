package register_handler

import (
	"bufio"
	"context"
	"log"
	"math/rand/v2"
	"net/http"
	"strconv"
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
	Events(fiberApp, mockSessionFinder, mockEventListener{})

	// Выбор порта
	port := strconv.Itoa(int(rand.Int32N(65535-49152) + 49152))
	// Запуск сервера
	go func() { assert.NoError(t, fiberApp.Listen(":"+port)) }()
	// Отложить остановку сервера
	defer assert.NoError(t, fiberApp.Shutdown())
	// Задержка чтобы сервер успел запуститься
	time.Sleep(time.Millisecond * 5)

	// Контекст для корректной отмены запроса
	ctx, cancel := context.WithCancel(context.Background())

	// Создать запрос
	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:"+port+"/events", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer 123")

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

func (m mockEventListener) AddListener(userID, sessionID uuid.UUID, f func(event any, err error)) (removeListener func(), err error) {
	go func() {
		for {
			f(rndDirtyEvent(), nil)
		}
	}()
	return func() {}, nil
}

func rndDirtyEvent() any {
	var e struct {
		Name string `fake:"{name}"`
		Age  int    `fake:"{number:0,100}"`
		Sex  string `fake:"{gender}"`
	}
	_ = gofakeit.Struct(&e)
	return e
}
