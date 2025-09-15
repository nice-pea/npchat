package jwtСreate

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_JWTC_Issue(t *testing.T) {
	t.Run("session может быть zero value", func(t *testing.T) {
		jwtC := Issuer{[]byte("secret")}

		session := sessionn.Session{}
		token, err := jwtC.Issue(session)
		require.NoError(t, err)
		assert.NotZero(t, token)
	})

	t.Run("jwt токены созданные с разными secret - неравны", func(t *testing.T) {
		var session = sessionn.Session{}

		jwtC := Issuer{[]byte("secret1")}
		t1, err := jwtC.Issue(session)
		require.NoError(t, err)

		jwtC = Issuer{[]byte("secret2")}
		t2, err := jwtC.Issue(session)
		require.NoError(t, err)

		assert.NotEqual(t, t1, t2)
	})
	t.Run("jwt токены созданные с zero secret - невалидны", func(t *testing.T) {
		var session = sessionn.Session{}

		jwtC := Issuer{}
		token, err := jwtC.Issue(session)
		require.Error(t, err)
		assert.Zero(t, token)
	})

	t.Run("jwt токены созданные с разными даными - неравны", func(t *testing.T) {
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

		jwtC := Issuer{[]byte("secret")}
		t1, err := jwtC.Issue(session1)
		require.NoError(t, err)
		t2, err := jwtC.Issue(session2)
		require.NoError(t, err)

		assert.NotEqual(t, t1, t2)
	})

	t.Run("Issue использует только ID и UserID, другие данные игнорируются", func(t *testing.T) {

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

		jwtC := Issuer{[]byte("secret")}
		t1, err := jwtC.Issue(session1)
		require.NoError(t, err)
		t2, err := jwtC.Issue(session2)
		require.NoError(t, err)

		assert.Equal(t, t1, t2)
	})
}
