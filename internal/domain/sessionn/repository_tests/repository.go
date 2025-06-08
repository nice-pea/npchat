package repository_tests

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
)

func TestRepository(t *testing.T, newRepository func() sessionn.Repository) {
	t.Run("List", func(t *testing.T) {
		t.Run("из пустого репозитория вернется пустой список", func(t *testing.T) {
			r := newRepository()
			users, err := r.List(sessionn.Filter{})
			assert.NoError(t, err)
			assert.Empty(t, users)
		})
	})
	t.Run("Upsert", func(t *testing.T) {
		t.Run("нельзя сохранять без ID", func(t *testing.T) {
			r := newRepository()
			err := r.Upsert(sessionn.Session{
				ID:   "",
				Name: "someName",
			})
			assert.Error(t, err)
		})

		t.Run("остальные поля, кроме ID могут быть пустыми", func(t *testing.T) {
			r := newRepository()
			err := r.Upsert(sessionn.Session{
				ID: uuid.NewString(),
			})
			assert.NoError(t, err)
		})

		t.Run("сохраненная сущность полностью соответствует сохраняемой", func(t *testing.T) {
			r := newRepository()
			// Создать
			session := rndUser(t)
			addRndBasicAuth(t, &session)
			addRndOpenAuth(t, &session)

			// Сохранить
			err := r.Upsert(session)
			require.NoError(t, err)

			// Прочитать из репозитория
			users, err := r.List(sessionn.Filter{})
			assert.NoError(t, err)
			require.Len(t, users, 1)
			assert.Equal(t, session, users[0])
		})

		t.Run("перезапись с новыми значениями по ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			// Несколько промежуточных состояний
			for range 33 {
				session := rndUser(t)
				session.ID = id
				upsertUser(t, r, session)
			}
			// Последнее сохраненное состояние
			expectedUser := rndUser(t)
			expectedUser.ID = id
			upsertUser(t, r, expectedUser)

			// Прочитать из репозитория
			users, err := r.List(sessionn.Filter{})
			assert.NoError(t, err)
			require.Len(t, users, 1)
			assert.Equal(t, expectedUser, users[0])
		})
	})
}

func rndUser(t *testing.T) sessionn.Session {
	u, err := sessionn.NewUser(gofakeit.Name(), gofakeit.Noun())
	require.NoError(t, err)
	return u
}

func addRndBasicAuth(t *testing.T, session *sessionn.Session) {
	ba, err := sessionn.NewBasicAuth(gofakeit.Noun(), "Passw0rd!")
	require.NoError(t, err)
	err = session.AddBasicAuth(ba)
	require.NoError(t, err)
}

func addRndOpenAuth(t *testing.T, session *sessionn.Session) {
	token, err := sessionn.NewOpenAuthToken(uuid.NewString(), "test", uuid.NewString(), time.Now().Add(1*time.Hour))
	require.NoError(t, err)
	openAuthUser, err := sessionn.NewOpenAuthUser(uuid.NewString(), "test", "", gofakeit.Noun(), "", token)
	require.NoError(t, err)
	err = session.AddOpenAuthUser(openAuthUser)
	require.NoError(t, err)
}

func upsertUser(t *testing.T, r sessionn.Repository, session sessionn.Session) sessionn.Session {
	err := r.Upsert(session)
	require.NoError(t, err)

	return session
}
