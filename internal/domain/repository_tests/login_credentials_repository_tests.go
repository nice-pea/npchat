package repository_tests

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func AuthnPasswordRepositoryTests(t *testing.T, newRepository func() domain.AuthnPasswordRepository) {
	t.Run("List", func(t *testing.T) {
		t.Run("без фильтра в пустом репозитории", func(t *testing.T) {
			r := newRepository()
			aps, err := r.List(domain.AuthnPasswordFilter{})
			assert.NoError(t, err)
			assert.Empty(t, aps)
		})
		t.Run("без фильтра только один AuthnPassword", func(t *testing.T) {
			r := newRepository()
			lc := domain.AuthnPassword{
				UserID:   uuid.NewString(),
				Login:    uuid.NewString(),
				Password: uuid.NewString(),
			}
			err := r.Save(lc)
			require.NoError(t, err)
			aps, err := r.List(domain.AuthnPasswordFilter{})
			assert.NoError(t, err)
			assert.Len(t, aps, 1)
		})
		t.Run("фильтр по userID", func(t *testing.T) {
			r := newRepository()
			lc := domain.AuthnPassword{
				UserID:   uuid.NewString(),
				Login:    uuid.NewString(),
				Password: uuid.NewString(),
			}
			err := r.Save(lc)
			require.NoError(t, err)
			aps, err := r.List(domain.AuthnPasswordFilter{UserID: lc.UserID})
			assert.NoError(t, err)
			if assert.Len(t, aps, 1) {
				assert.Equal(t, lc.UserID, aps[0].UserID)
			}
		})
		t.Run("фильтр по login и password", func(t *testing.T) {
			r := newRepository()
			var desired domain.AuthnPassword
			// Создаем  AuthnPassword
			for range 10 {
				lc := domain.AuthnPassword{
					UserID:   uuid.NewString(),
					Login:    uuid.NewString(),
					Password: uuid.NewString(),
				}
				err := r.Save(lc)
				require.NoError(t, err)
				// Сохранить созданную запись
				desired = lc
			}
			aps, err := r.List(domain.AuthnPasswordFilter{
				Login:    desired.Login,
				Password: desired.Password,
			})
			assert.NoError(t, err)
			require.Len(t, aps, 1)
			assert.Equal(t, desired, aps[0])
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("нельзя сохранять AuthnPassword без userID", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.AuthnPassword{
				UserID:   "",
				Login:    uuid.NewString(),
				Password: uuid.NewString(),
			})
			assert.Error(t, err)
		})
		t.Run("сохраненную AuthnPassword можно прочитать из репозитория", func(t *testing.T) {
			r := newRepository()
			session := domain.AuthnPassword{
				UserID:   uuid.NewString(),
				Login:    uuid.NewString(),
				Password: uuid.NewString(),
			}
			err := r.Save(session)
			require.NoError(t, err)
			creds, err := r.List(domain.AuthnPasswordFilter{})
			assert.NoError(t, err)
			require.Len(t, creds, 1)
			assert.Equal(t, session, creds[0])
		})
	})
}
