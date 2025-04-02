package repository_tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/stretchr/testify/assert"
)

func UsersRepositoryTests(t *testing.T, newRepository func() domain.UsersRepository) {
	t.Run("List", func(t *testing.T) {
		t.Run("без фильтра в пустом репозитории", func(t *testing.T) {
			r := newRepository()
			users, err := r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			assert.Len(t, users, 0)
		})
		t.Run("без фильтра только один чат", func(t *testing.T) {
			r := newRepository()
			user := domain.User{
				ID: uuid.NewString(),
			}
			err := r.Save(user)
			assert.NoError(t, err)
			users, err := r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			assert.Len(t, users, 1)
		})
		t.Run("без фильтра отсутствует после удаления", func(t *testing.T) {
			r := newRepository()
			user := domain.User{
				ID: uuid.NewString(),
			}
			err := r.Save(user)
			assert.NoError(t, err)
			users, err := r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			assert.Len(t, users, 1)
			err = r.Delete(user.ID)
			assert.NoError(t, err)
			users, err = r.List(domain.UsersFilter{})
			assert.NoError(t, err)
			assert.Len(t, users, 0)
		})
		t.Run("с фильтром по id", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			assert.NoError(t, errors.Join(
				r.Save(domain.User{ID: id}),
				r.Save(domain.User{ID: uuid.NewString()}),
				r.Save(domain.User{ID: uuid.NewString()}),
			))
			users, err := r.List(domain.UsersFilter{ID: id})
			assert.NoError(t, err)
			if assert.Len(t, users, 1) {
				assert.Equal(t, id, users[0].ID)
			}
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("без ID", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.User{
				ID: "",
			})
			assert.Error(t, err)
		})
		t.Run("без ошибок", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.User{
				ID: uuid.NewString(),
			})
			assert.NoError(t, err)
		})
		t.Run("дубль ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.User{
				ID: id,
			})
			assert.NoError(t, err)
			err = r.Save(domain.User{
				ID: id,
			})
			assert.NoError(t, err) // с чем связанно, тут не возвращается error?
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("с пустым id", func(t *testing.T) {
			r := newRepository()
			err := r.Delete("")
			assert.Error(t, err)
		})
		t.Run("несуществующий id", func(t *testing.T) {
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
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
		t.Run("дважды удаленный", func(t *testing.T) {
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
