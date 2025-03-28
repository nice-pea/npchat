package repository_tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func MembersRepositoryTests(t *testing.T, newRepository func() domain.MembersRepository) {
	t.Helper()
	t.Run("List", func(t *testing.T) {
		t.Run("без фильтра в пустом репозитории", func(t *testing.T) {
			r := newRepository()
			members, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			assert.Len(t, members, 0)
		})
		t.Run("без фильтра только один участник", func(t *testing.T) {
			r := newRepository()
			member := domain.Member{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
			}
			err := r.Save(member)
			assert.NoError(t, err)
			members, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			assert.Len(t, members, 1)
		})
		t.Run("без фильтра отсутствует после удаления", func(t *testing.T) {
			r := newRepository()
			member := domain.Member{
				ID:     uuid.NewString(),
				ChatID: "name",
			}
			err := r.Save(member)
			assert.NoError(t, err)
			members, err := r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			assert.Len(t, members, 1)
			err = r.Delete(member.ID)
			assert.NoError(t, err)
			members, err = r.List(domain.MembersFilter{})
			assert.NoError(t, err)
			assert.Len(t, members, 0)
		})
		t.Run("с фильтром по id", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			assert.NoError(t, errors.Join(
				r.Save(domain.Member{ID: id, ChatID: uuid.NewString()}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString()}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString()}),
			))
			members, err := r.List(domain.MembersFilter{ID: id})
			assert.NoError(t, err)
			assert.Len(t, members, 1)
		})
		t.Run("с фильтром по chat id", func(t *testing.T) {
			r := newRepository()
			chatID := uuid.NewString()
			assert.NoError(t, errors.Join(
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: chatID}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: chatID}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString()}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString()}),
			))
			members, err := r.List(domain.MembersFilter{ChatID: chatID})
			assert.NoError(t, err)
			assert.Len(t, members, 2)
		})
		t.Run("с фильтром по user id", func(t *testing.T) {
			r := newRepository()
			userID := uuid.NewString()
			assert.NoError(t, errors.Join(
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: userID}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: userID}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString()}),
				r.Save(domain.Member{ID: uuid.NewString(), ChatID: uuid.NewString(), UserID: uuid.NewString()}),
			))
			members, err := r.List(domain.MembersFilter{UserID: userID})
			assert.NoError(t, err)
			assert.Len(t, members, 2)
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("чат без id", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Member{
				ID:     "",
				ChatID: "name",
			})
			assert.Error(t, err)
		})
		t.Run("чат без name", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Member{
				ID:     uuid.NewString(),
				ChatID: "",
			})
			assert.NoError(t, err)
		})
		t.Run("без ошибок", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Member{
				ID:     uuid.NewString(),
				ChatID: "name",
			})
			assert.NoError(t, err)
		})
		t.Run("дубль ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Member{
				ID:     id,
				ChatID: "name",
			})
			assert.NoError(t, err)
			err = r.Save(domain.Member{
				ID:     id,
				ChatID: "name1",
			})
			assert.Error(t, err)
		})
		t.Run("дубль name", func(t *testing.T) {
			r := newRepository()
			name := "name"
			err := r.Save(domain.Member{
				ID:     uuid.NewString(),
				ChatID: name,
			})
			assert.NoError(t, err)
			err = r.Save(domain.Member{
				ID:     uuid.NewString(),
				ChatID: name,
			})
			assert.NoError(t, err)
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
			err := r.Save(domain.Member{
				ID:     id,
				ChatID: "name",
			})
			err = r.Delete(id)
			assert.NoError(t, err)
		})
		t.Run("дважды удаленный", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Member{
				ID:     id,
				ChatID: "name",
			})
			err = r.Delete(id)
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
	})
}
