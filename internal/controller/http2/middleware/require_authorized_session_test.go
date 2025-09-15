package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
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
	t.Run("sessions", func(t *testing.T) {
		t.Run("сохраненную SessionID и UserID можно прочитать", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{
				FindSessionsFunc: func(in findSession.In) (findSession.Out, error) {
					return findSession.Out{Sessions: []sessionn.Session{mockSession}}, nil
				},
			}

			fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
			fiberApp.Get(
				"/", RequireAuthorizedSession(uc, nil),
				func(ctx *fiber.Ctx) error {
					sessionId, ok := ctx.Locals(CtxKeySessionID).(uuid.UUID)
					require.True(t, ok)
					userId, ok := ctx.Locals(CtxKeyUserID).(uuid.UUID)
					require.True(t, ok)

					require.Equal(t, mockSession.ID, sessionId)
					require.Equal(t, mockSession.UserID, userId)
					return nil
				})

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "SessionToken "+mockSession.AccessToken.Token)

			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})

		t.Run("не передан или с ошибкой в ключе заголовока - вернет StatusUnauthorized (401)", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}
			fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
			fiberApp.Get("/", RequireAuthorizedSession(uc, nil))

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorizati1on", "SessionToken "+mockSession.AccessToken.Token)

			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run("FindSessions вернет ошибку StatusInternalServerError (500)", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{
				FindSessionsFunc: func(in findSession.In) (findSession.Out, error) {
					return findSession.Out{}, errors.New("БД отключена")
				},
			}

			fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
			fiberApp.Get("/", RequireAuthorizedSession(uc, nil))

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "SessionToken "+mockSession.AccessToken.Token)

			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		})
		t.Run("сессия не найдена - 401", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{
				FindSessionsFunc: func(in findSession.In) (findSession.Out, error) {
					return findSession.Out{}, nil
				},
			}

			fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
			fiberApp.Get("/", RequireAuthorizedSession(uc, nil))

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "SessionToken "+mockSession.AccessToken.Token)

			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run("неправильный формат заголовка - 401", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}
			fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
			fiberApp.Get("/", RequireAuthorizedSession(uc, nil))

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "InvalidFormat "+mockSession.AccessToken.Token)

			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run("пустой токен после SessionToken - 401", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}
			fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
			fiberApp.Get("/", RequireAuthorizedSession(uc, nil))

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "SessionToken ")

			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
	})
	t.Run("JWT", func(t *testing.T) {
		t.Run("валидный JWT", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}
			mtm := mockJWTParser{}

			fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
			fiberApp.Get(
				"/", RequireAuthorizedSession(uc, mtm),
				func(ctx *fiber.Ctx) error {
					userid, ok := ctx.Locals("UserID").(string)
					require.True(t, ok)
					require.Equal(t, mockParseJWT.UserID, userid)
					sessionid, ok := ctx.Locals("SessionID").(string)
					require.True(t, ok)
					require.Equal(t, mockParseJWT.SessionID, sessionid)
					return nil
				})

			jwTocken := "31132"

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer "+jwTocken)

			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
		t.Run("истекший JWT вернет StatusUnauthorized", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}
			mtm := mockJWTParser{
				ParseFunc: func(token string) (OutJWT, error) {
					return OutJWT{}, errors.New("JWT истекший")
				},
			}

			fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
			fiberApp.Get(
				"/", RequireAuthorizedSession(uc, mtm),
				func(ctx *fiber.Ctx) error {
					assert.Fail(t, "unreachable code")
					return nil
				})

			jwTocken := "123.456.789"

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer "+jwTocken)

			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})
		t.Run("пустое значение jwt токена", func(t *testing.T) {
			uc := mockUsecasesForRequireAuthorizedSession{}
			mtm := mockJWTParser{}

			fiberApp := fiber.New(fiber.Config{DisableStartupMessage: true})
			fiberApp.Get(
				"/", RequireAuthorizedSession(uc, mtm),
				func(ctx *fiber.Ctx) error {
					assert.Fail(t, "unreachable code")
					return nil
				})

			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("Authorization", "Bearer ")

			resp, err := fiberApp.Test(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
		})

	})
}

type mockUsecasesForRequireAuthorizedSession struct {
	FindSessionsFunc func(findSession.In) (findSession.Out, error)
}

var mockSession = sessionn.Session{
	ID:     uuid.New(),
	Status: sessionn.StatusNew,
	UserID: uuid.New(),
	Name:   "name",
	AccessToken: sessionn.Token{
		Token: "asdasda",
	},
}

func (m mockUsecasesForRequireAuthorizedSession) FindSessions(in findSession.In) (findSession.Out, error) {
	if m.FindSessionsFunc != nil {
		return m.FindSessionsFunc(in)
	}
	return findSession.Out{Sessions: []sessionn.Session{mockSession}}, nil
}

type mockJWTParser struct {
	ParseFunc func(token string) (OutJWT, error)
}

var mockParseJWT = OutJWT{
	UserID:    "1234",
	SessionID: "5678",
}

func (m mockJWTParser) Parse(token string) (OutJWT, error) {
	if m.ParseFunc != nil {
		return m.ParseFunc(token)
	}
	return mockParseJWT, nil
}
