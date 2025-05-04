package repository_tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func LoginCredentialsRepositoryTests(t *testing.T, newRepository func() domain.LoginCredentialsRepository) {
	t.Run("List", func(t *testing.T) {
		t.Run("без фильтра в пустом репозитории", func(t *testing.T) {
			r := newRepository()
			creds, err := r.List(domain.LoginCredentialsFilter{})
			assert.NoError(t, err)
			assert.Empty(t, creds)
		})
		t.Run("без фильтра только одни credentials", func(t *testing.T) {
			r := newRepository()
			lc := domain.LoginCredentials{
				UserID:   uuid.NewString(),
				Login:    uuid.NewString(),
				Password: uuid.NewString(),
			}
			err := r.Save(lc)
			require.NoError(t, err)
			creds, err := r.List(domain.LoginCredentialsFilter{})
			assert.NoError(t, err)
			assert.Len(t, creds, 1)
		})
		t.Run("фильтр по userID", func(t *testing.T) {
			r := newRepository()
			lc := domain.LoginCredentials{
				UserID:   uuid.NewString(),
				Login:    uuid.NewString(),
				Password: uuid.NewString(),
			}
			err := r.Save(lc)
			require.NoError(t, err)
			creds, err := r.List(domain.LoginCredentialsFilter{UserID: lc.UserID})
			assert.NoError(t, err)
			if assert.Len(t, creds, 1) {
				assert.Equal(t, lc.UserID, creds[0].UserID)
			}
		})
		t.Run("фильтр по login и password", func(t *testing.T) {
			r := newRepository()
			var desired domain.LoginCredentials
			// Создаем  credentials
			for range 10 {
				lc := domain.LoginCredentials{
					UserID:   uuid.NewString(),
					Login:    uuid.NewString(),
					Password: uuid.NewString(),
				}
				err := r.Save(lc)
				require.NoError(t, err)
				// Сохранить созданную запись
				desired = lc
			}
			creds, err := r.List(domain.LoginCredentialsFilter{
				Login:    desired.Login,
				Password: desired.Password,
			})
			assert.NoError(t, err)
			require.Len(t, creds, 1)
			assert.Equal(t, desired, creds[0])
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("нельзя сохранять credentials без userID", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.LoginCredentials{
				UserID:   "",
				Login:    uuid.NewString(),
				Password: uuid.NewString(),
			})
			assert.Error(t, err)
		})
		t.Run("сохраненную credentials можно прочитать из репозитория", func(t *testing.T) {
			r := newRepository()
			session := domain.LoginCredentials{
				UserID:   uuid.NewString(),
				Login:    uuid.NewString(),
				Password: uuid.NewString(),
			}
			err := r.Save(session)
			require.NoError(t, err)
			creds, err := r.List(domain.LoginCredentialsFilter{})
			assert.NoError(t, err)
			require.Len(t, creds, 1)
			assert.Equal(t, session, creds[0])
		})
	})
}
