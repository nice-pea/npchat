package middleware

import (
	"errors"
	"log"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// middleware_mock "github.com/nice-pea/npchat/internal/controller/http2/middleware/mocks"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	findSession "github.com/nice-pea/npchat/internal/usecases/sessions/find_session"
)

func Test_RequireAuthorizedSession(t *testing.T) {
	log.Print("тест RequireAuthorizedSession")
	t.Run("sessions", func(t *testing.T) {
		t.Run("сохраненную сессию можно прочитать", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{
				FindSessionsFunc: func(in findSession.In) (findSession.Out, error) {
					return findSession.Out{Sessions: []sessionn.Session{mockSession}}, nil
				},
			}

			server := fiber.New(fiber.Config{DisableStartupMessage: true})
			server.Get(
				"/", RequireAuthorizedSession(uc),
				func(ctx *fiber.Ctx) error {
					session, ok := ctx.Locals(CtxKeyUserSession).(sessionn.Session)
					require.True(t, ok)
					require.Equal(t, mockSession, session)
					return nil
				})

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+mockSession.AccessToken.Token)

			resp, err := server.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("не передан или с ошибкой в ключе заголовока - вернет StatusUnauthorized (401)", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}
			server := fiber.New(fiber.Config{DisableStartupMessage: true})
			server.Get("/", RequireAuthorizedSession(uc))

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorizati1on", "Bearer "+mockSession.AccessToken.Token)

			resp, err := server.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run("FindSessions вернет ошибку StatusInternalServerError (500)", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{
				FindSessionsFunc: func(in findSession.In) (findSession.Out, error) {
					return findSession.Out{}, errors.New("БД отключена")
				},
			}

			server := fiber.New(fiber.Config{DisableStartupMessage: true})
			server.Get("/", RequireAuthorizedSession(uc))

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+mockSession.AccessToken.Token)

			resp, err := server.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		})
		t.Run("сессия не найдена - 401", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{
				FindSessionsFunc: func(in findSession.In) (findSession.Out, error) {
					return findSession.Out{}, nil
				},
			}

			server := fiber.New(fiber.Config{DisableStartupMessage: true})
			server.Get("/", RequireAuthorizedSession(uc))

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer "+mockSession.AccessToken.Token)

			resp, err := server.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run("неправильный формат заголовка - 401", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}
			server := fiber.New(fiber.Config{DisableStartupMessage: true})
			server.Get("/", RequireAuthorizedSession(uc))

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "InvalidFormat "+mockSession.AccessToken.Token)

			resp, err := server.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run("пустой токен после Bearer - 401", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}
			server := fiber.New(fiber.Config{DisableStartupMessage: true})
			server.Get("/", RequireAuthorizedSession(uc))

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "Bearer ")

			resp, err := server.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
	t.Run("JWT", func(t *testing.T) {
		t.Run("валидный JWT", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}

			server := fiber.New(fiber.Config{DisableStartupMessage: true})
			server.Get(
				"/", RequireAuthorizedSession(uc),
				func(ctx *fiber.Ctx) error {
					session, ok := ctx.Locals(CtxKeyUserSession).(sessionn.Session)
					require.True(t, ok)
					require.Equal(t, mockSession, session)
					return nil
				})

			jwTocken := ""

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "JWT "+jwTocken)

			resp, err := server.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
		t.Run("истекший JWT", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}

			server := fiber.New(fiber.Config{DisableStartupMessage: true})
			server.Get(
				"/", RequireAuthorizedSession(uc),
				func(ctx *fiber.Ctx) error {
					session, ok := ctx.Locals(CtxKeyUserSession).(sessionn.Session)
					require.True(t, ok)
					require.Equal(t, mockSession, session)
					return nil
				})

			jwTocken := ""

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "JWT "+jwTocken)

			resp, err := server.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run("не верной подписи JWT", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}

			server := fiber.New(fiber.Config{DisableStartupMessage: true})
			server.Get(
				"/", RequireAuthorizedSession(uc),
				func(ctx *fiber.Ctx) error {
					session, ok := ctx.Locals(CtxKeyUserSession).(sessionn.Session)
					require.True(t, ok)
					require.Equal(t, mockSession, session)
					return nil
				})

			jwTocken := ""

			req, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", "JWT "+jwTocken)

			resp, err := server.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		

	})
}

type mockUsecasesForRequireAuthorizedSession struct {
	FindSessionsFunc func(findSession.In) (findSession.Out, error)
}

var (
	mockSession = sessionn.Session{
		ID:     uuid.New(),
		Status: sessionn.StatusNew,
		UserID: uuid.New(),
		Name:   "name",
		AccessToken: sessionn.Token{
			Token: "asdasda",
		},
	}
)

func (m mockUsecasesForRequireAuthorizedSession) FindSessions(in findSession.In) (findSession.Out, error) {
	if m.FindSessionsFunc != nil {
		return m.FindSessionsFunc(in)
	}
	return findSession.Out{Sessions: []sessionn.Session{mockSession}}, nil
}

func CreateJWT() string {
	return ""
}
