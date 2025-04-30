package repository_tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func SessionsRepositoryTests(t *testing.T, newRepository func() domain.SessionsRepository) {
	t.Run("List", func(t *testing.T) {
		t.Run("без фильтра в пустом репозитории", func(t *testing.T) {
			r := newRepository()
			sessions, err := r.List(domain.SessionsFilter{})
			assert.NoError(t, err)
			assert.Empty(t, sessions)
		})
		t.Run("без фильтра только одна сессия", func(t *testing.T) {
			r := newRepository()
			session := domain.Session{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				Token:  uuid.NewString(),
				Status: domain.SessionStatusNew,
			}
			err := r.Save(session)
			assert.NoError(t, err)
			sessions, err := r.List(domain.SessionsFilter{})
			assert.NoError(t, err)
			assert.Len(t, sessions, 1)
		})
		t.Run("фильтр по токену", func(t *testing.T) {
			r := newRepository()
			token := "test-token"
			session := domain.Session{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				Token:  token,
				Status: domain.SessionStatusNew,
			}
			err := r.Save(session)
			assert.NoError(t, err)
			sessions, err := r.List(domain.SessionsFilter{Token: token})
			assert.NoError(t, err)
			if assert.Len(t, sessions, 1) {
				assert.Equal(t, token, sessions[0].Token)
			}
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("нельзя сохранять сессию без ID", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Session{
				ID: "",
			})
			assert.Error(t, err)
		})
		t.Run("сохраненную сессию можно прочитать из репозитория", func(t *testing.T) {
			r := newRepository()
			session := domain.Session{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				Token:  uuid.NewString(),
				Status: domain.SessionStatusNew,
			}
			err := r.Save(session)
			assert.NoError(t, err)
			sessions, err := r.List(domain.SessionsFilter{})
			assert.NoError(t, err)
			require.Len(t, sessions, 1)
			assert.Equal(t, session, sessions[0])
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("ID обязательно должен быть передан", func(t *testing.T) {
			r := newRepository()
			err := r.Delete("")
			assert.Error(t, err)
		})
		t.Run("удаление несуществующей записи не вернет ошибку", func(t *testing.T) {
			r := newRepository()
			err := r.Delete(uuid.NewString())
			assert.NoError(t, err)
		})
		t.Run("без ошибок", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Session{
				ID:     id,
				UserID: uuid.NewString(),
				Token:  uuid.NewString(),
				Status: domain.SessionStatusNew,
			})
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
		t.Run("можно несколько раз удалять", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Session{
				ID:     id,
				UserID: uuid.NewString(),
				Token:  uuid.NewString(),
				Status: domain.SessionStatusNew,
			})
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
	})
}
