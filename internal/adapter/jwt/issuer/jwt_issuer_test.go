package jwtIssuer

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	jwt2 "github.com/nice-pea/npchat/internal/adapter/jwt"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

// Test_Issuer_Issue набор тестов для проверки функции Issue генератора JWT-токенов
func Test_Issuer_Issue(t *testing.T) {
	t.Run("session может быть zero value", func(t *testing.T) {
		issuer := Issuer{jwt2.Config{SecretKey: "secret"}}

		session := sessionn.Session{}
		token, err := issuer.Issue(session)
		require.NoError(t, err)
		assert.NotZero(t, token)
	})

	t.Run("jwt токены созданные с разными secret - неравны", func(t *testing.T) {
		var session = sessionn.Session{}

		issuer := Issuer{jwt2.Config{SecretKey: "secret1"}}
		t1, err := issuer.Issue(session)
		require.NoError(t, err)

		issuer = Issuer{jwt2.Config{SecretKey: "secret2"}}
		t2, err := issuer.Issue(session)
		require.NoError(t, err)

		assert.NotEqual(t, t1, t2)
	})

	t.Run("jwt токены созданные с zero secret - невалидны", func(t *testing.T) {
		var session = sessionn.Session{}

		issuer := Issuer{}
		token, err := issuer.Issue(session)
		require.Error(t, err)
		assert.Zero(t, token)
	})

	t.Run("jwt токены созданные с разными данными - неравны", func(t *testing.T) {
		var (
			session1 = sessionn.Session{
				ID:     uuid.New(),
				UserID: uuid.New(),
			}
			session2 = sessionn.Session{
				ID:     uuid.New(),
				UserID: uuid.New(),
			}
		)

		issuer := Issuer{jwt2.Config{SecretKey: "secret"}}
		t1, err := issuer.Issue(session1)
		require.NoError(t, err)
		t2, err := issuer.Issue(session2)
		require.NoError(t, err)

		assert.NotEqual(t, t1, t2)
	})

	t.Run("Issue использует только ID и UserID, другие данные игнорируются", func(t *testing.T) {
		// в этом тесте предпологается что jwt создаются в одно время
		var (
			ID     = uuid.New()
			UserID = uuid.New()

			session1 = sessionn.Session{
				ID:     ID,
				UserID: UserID,
				Name:   "session1",
				Status: "session1",
				AccessToken: sessionn.Token{
					Token:  "session1",
					Expiry: time.Time{},
				},
				RefreshToken: sessionn.Token{
					Token:  "session1",
					Expiry: time.Time{},
				},
			}
			session2 = sessionn.Session{
				ID:     ID,
				UserID: UserID,
				Name:   "session2",
				Status: "session2",
				AccessToken: sessionn.Token{
					Token:  "session2",
					Expiry: time.Time{},
				},
				RefreshToken: sessionn.Token{
					Token:  "session2",
					Expiry: time.Time{},
				},
			}
		)

		issuer := Issuer{jwt2.Config{SecretKey: "secret"}}
		t1, err := issuer.Issue(session1)
		require.NoError(t, err)
		t2, err := issuer.Issue(session2)
		require.NoError(t, err)

		assert.Equal(t, t1, t2)
	})

	t.Run("Также Issue использует ExpiresAt, в разное время созданные токены не равны", func(t *testing.T) {

		var session = sessionn.Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
		}

		issuer := Issuer{jwt2.Config{SecretKey: "secret"}}
		t1, err := issuer.Issue(session)
		require.NoError(t, err)

		time.Sleep(time.Second)

		t2, err := issuer.Issue(session)
		require.NoError(t, err)

		assert.NotEqual(t, t1, t2)
	})

	t.Run("Если secret пустой - Issue будет возвращать ошибку", func(t *testing.T) {

		var session = sessionn.Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
		}

		issuer := Issuer{jwt2.Config{SecretKey: ""}}
		token, err := issuer.Issue(session)
		require.Error(t, err)

		assert.Zero(t, token)
	})

	t.Run("в jwt запсывается дата создания токена IssuedAt", func(t *testing.T) {
		secret := "qwerty"
		var session = sessionn.Session{
			ID:     uuid.New(),
			UserID: uuid.New(),
		}
		issuer := Issuer{jwt2.Config{SecretKey: secret}}

		issuedAtStart := time.Now()
		token, err := issuer.Issue(session)
		issuedAtEnd := time.Now()

		require.NoError(t, err)
		claims := parse(t, token, []byte(secret))
		issuedAt := claims.IssuedAt
		require.Equal(t, session.ID, claims.SessionID)
		require.Equal(t, session.UserID, claims.UserID)

		require.GreaterOrEqual(t, issuedAt.Unix(), issuedAtStart.Unix())
		require.LessOrEqual(t, issuedAt.Unix(), issuedAtEnd.Unix())
	})
}

// claims структура для хранения декодированных данных из JWT токена
type claims struct {
	UserID    uuid.UUID
	SessionID uuid.UUID
	jwt.RegisteredClaims
}

// parse вспомогательная функция для декодирования и проверки JWT токена
func parse(t *testing.T, token string, secret []byte) claims {
	t.Helper()
	verifier, err := jwt.NewVerifierHS(jwt.HS256, secret)
	require.NoError(t, err)

	tokenBytes := []byte(token)
	newToken, err := jwt.Parse(tokenBytes, verifier)
	require.NoError(t, err)

	err = verifier.Verify(newToken)
	require.NoError(t, err)

	var newClaims claims
	errClaims := json.Unmarshal(newToken.Claims(), &newClaims)
	require.NoError(t, errClaims)

	require.True(t, newClaims.IsValidAt(time.Now()))

	return newClaims
}
