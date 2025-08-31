package register_handler

import (
	"bufio"
	"context"
	"fmt"
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
		IdleTimeout:           time.Second,
		DisableKeepalive:      false,
		DisableStartupMessage: true,
		StreamRequestBody:     true,
	})
	Events(server, mockSessionFinder{}, mockEventListener{})
	go func() { assert.NoError(t, server.Listen("localhost:8418")) }()
	time.Sleep(time.Millisecond * 100)

	// Создаем HTTP-запрос
	resp, err := http.Get("http://localhost:8418/events")
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	// Читаем данные из потока
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) > 0 {
			fmt.Println("Получено событие:", line)
		}
	}
	require.NoError(t, scanner.Err())
	assert.NoError(t, server.Shutdown())
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

func (m mockEventListener) Listen(ctx context.Context, userID uuid.UUID, sessionID uuid.UUID, f func(event any)) error {

	f(map[string]any{
		"name": gofakeit.Name(),
		"id":   gofakeit.UUID(),
		"age":  gofakeit.Number(0, 100),
		"sex":  gofakeit.Gender(),
	})

	select {
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}
