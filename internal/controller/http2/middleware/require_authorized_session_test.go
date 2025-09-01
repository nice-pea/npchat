package middleware

import (
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	// middleware_mock "github.com/nice-pea/npchat/internal/controller/http2/middleware/mocks"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	findSession "github.com/nice-pea/npchat/internal/usecases/sessions/find_session"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RequireAuthorizedSession(t *testing.T) {
	log.Print("тест RequireAuthorizedSession")
	t.Run("сохраненную сессию можно прочитать", func(t *testing.T) {
		uc := mockUsecasesForRequireAuthorizedSession{}

		// var wg sync.WaitGroup

		// wg.Add(1)
		server := fiber.New(fiber.Config{DisableStartupMessage: true})
		server.Get(
			"/", RequireAuthorizedSession(uc),
			func(ctx *fiber.Ctx) error {
				// defer wg.Done()
				sessionFromLocals := ctx.Locals(CtxKeyUserSession, sessionn.Session{})
				require.IsType(t, sessionn.Session{}, sessionFromLocals)
				log.Print("значение из locals это сессия")
				session := sessionFromLocals.(sessionn.Session)
				require.Equal(t, mockSession, session)
				log.Print("сессии одинаковые")
				return nil
			})

		go func() { assert.NoError(t, server.Listen("localhost:8419")) }()
		time.Sleep(time.Millisecond * 10)

		req, err := http.NewRequest("GET", "http://localhost:8419/", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", "Bearer "+mockSession.AccessToken.Token)

		_, err = http.DefaultClient.Do(req)
		require.NoError(t, err)

		// wg.Wait()

		assert.NoError(t, server.Shutdown())
	})
	// t.Run("сохраненную сессию можно прочитать", func(t *testing.T) {
	// 	knownToken := mockSession.AccessToken.Token
	// 	uc := mockUsecasesForRequireAuthorizedSession{}
	// 	mw := RequireAuthorizedSession(&uc)

	// 	ctx := new(fiber.Ctx)
	// 	ctx.Set("Authorization", "Bearer "+knownToken)
	// 	mw(ctx)

	// 	sessionFromLocals := ctx.Locals(CtxKeyUserSession, sessionn.Session{})
	// 	require.IsType(t, sessionn.Session{}, sessionFromLocals)
	// 	session := sessionFromLocals.(sessionn.Session)
	// 	assert.Equal(t, mockSession, session)
	// })
}

type mockUsecasesForRequireAuthorizedSession struct{}

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
	return findSession.Out{
		Sessions: []sessionn.Session{mockSession},
	}, nil
}
