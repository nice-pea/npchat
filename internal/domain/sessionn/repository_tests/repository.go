package repository_tests

import (
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nice-pea/npchat/internal/common"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
)

func TestRepository(t *testing.T, newRepository func() sessionn.Repository) {
	t.Run("List", func(t *testing.T) {
		t.Run("из пустого репозитория вернется пустой список", func(t *testing.T) {
			r := newRepository()
			sessions, err := r.List(sessionn.Filter{})
			assert.NoError(t, err)
			assert.Empty(t, sessions)
		})

		t.Run("без фильтра из репозитория вернутся все сохраненные элементы", func(t *testing.T) {
			r := newRepository()
			sessions := upsertRndSessions(t, r, 10)
			fromRepo, err := r.List(sessionn.Filter{})
			assert.NoError(t, err)
			assert.Len(t, fromRepo, len(sessions))
		})

		t.Run("с фильтром по AccessToken вернутся, имеющие такой токен доступа", func(t *testing.T) {
			r := newRepository()
			// Создать много
			sessions := upsertRndSessions(t, r, 10)
			// Определить случайны искомый
			expected := common.RndElem(sessions)
			// Получить список
			fromRepo, err := r.List(sessionn.Filter{
				AccessToken: expected.AccessToken.Token,
			})
			// Сравнить ожидания и результат
			assert.NoError(t, err)
			require.Len(t, fromRepo, 1)
			assert.Equal(t, expected, fromRepo[0])
		})
	})

	t.Run("Upsert", func(t *testing.T) {
		t.Run("нельзя сохранять без ID", func(t *testing.T) {
			r := newRepository()
			err := r.Upsert(sessionn.Session{
				ID:   uuid.Nil,
				Name: "someName",
			})
			assert.Error(t, err)
		})

		t.Run("остальные поля, кроме ID могут быть пустыми", func(t *testing.T) {
			r := newRepository()
			err := r.Upsert(sessionn.Session{
				ID: uuid.New(),
			})
			assert.NoError(t, err)
		})

		t.Run("сохраненная сущность полностью соответствует сохраняемой", func(t *testing.T) {
			r := newRepository()
			// Создать и Сохранить
			session := upsertSession(t, r, rndSession(t))

			// Прочитать из репозитория
			sessions, err := r.List(sessionn.Filter{})
			assert.NoError(t, err)
			require.Len(t, sessions, 1)
			assert.Equal(t, session, sessions[0])
		})

		t.Run("перезапись с новыми значениями по ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.New()
			// Несколько промежуточных состояний
			for range 33 {
				session := rndSession(t)
				session.ID = id
				upsertSession(t, r, session)
			}
			// Последнее сохраненное состояние
			expected := rndSession(t)
			expected.ID = id
			upsertSession(t, r, expected)

			// Прочитать из репозитория
			sessions, err := r.List(sessionn.Filter{})
			assert.NoError(t, err)
			require.Len(t, sessions, 1)
			assert.Equal(t, expected, sessions[0])
		})
	})
}

func rndSession(t *testing.T) sessionn.Session {
	session, err := sessionn.NewSession(uuid.New(), gofakeit.UserAgent(), common.RndElem(sessionn.Statuses()))
	require.NoError(t, err)

	return session
}

func upsertRndSessions(t *testing.T, r sessionn.Repository, count int) []sessionn.Session {
	ss := make([]sessionn.Session, count)
	for i := range ss {
		ss[i] = rndSession(t)
		upsertSession(t, r, ss[i])
	}

	return ss
}

func upsertSession(t *testing.T, r sessionn.Repository, session sessionn.Session) sessionn.Session {
	err := r.Upsert(session)
	require.NoError(t, err)

	return session
}
