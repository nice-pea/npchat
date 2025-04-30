package repository_tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func UsersRepositoryTests(t *testing.T, newRepository func() domain.UsersRepository) {
	t.Run("List", func(t *testing.T) {
		t.Run("из пустого репозитория вернутся пустой список", func(t *testing.T) {
			r := newRepository()
			users, err := r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			assert.Empty(t, users)
		})
		t.Run("без фильтра вернутся все пользователи", func(t *testing.T) {
			r := newRepository()
			const usersCount = 11
			for range usersCount {
				require.NoError(t, r.Save(domain.User{
					ID: uuid.NewString(),
				}))
			}
			users, err := r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			assert.Len(t, users, usersCount)
		})
		t.Run("без фильтра отсутствует после удаления", func(t *testing.T) {
			r := newRepository()
			user := domain.User{
				ID: uuid.NewString(),
			}
			err := r.Save(user)
			require.NoError(t, err)
			users, err := r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			require.Len(t, users, 1)
			err = r.Delete(user.ID)
			require.NoError(t, err)
			users, err = r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			assert.Empty(t, users)
		})
		t.Run("с фильтром по id", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			require.NoError(t, errors.Join(
				r.Save(domain.User{ID: id}),
				r.Save(domain.User{ID: uuid.NewString()}),
				r.Save(domain.User{ID: uuid.NewString()}),
			))
			users, err := r.List(domain.UsersFilter{ID: id})
			assert.NoError(t, err)
			require.Len(t, users, 1)
			assert.Equal(t, id, users[0].ID)
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("ID обязательно должен быть передан", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.User{
				ID: "",
			})
			assert.Error(t, err)
		})
		t.Run("сохраненного пользователя можно прочитать из репозитория", func(t *testing.T) {
			r := newRepository()
			user := domain.User{
				ID: uuid.NewString(),
			}
			err := r.Save(user)
			assert.NoError(t, err)
			users, err := r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			require.Len(t, users, 1)
			assert.Equal(t, user, users[0])
		})
		t.Run("можно несколько раз сохранять с одним ID", func(t *testing.T) {
			r := newRepository()
			user := domain.User{ID: uuid.NewString()}
			for range 10 {
				require.NoError(t, r.Save(user))
			}
			users, err := r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			assert.Len(t, users, 1)
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
			err := r.Save(domain.User{
				ID: id,
			})
			require.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
		t.Run("можно несколько раз удалять", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.User{
				ID: id,
			})
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
	})
}
