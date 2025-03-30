package repository_tests

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func assertEqualChats(t *testing.T, expected, actual domain.Chat) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ChiefUserID, actual.ChiefUserID)
}

func ChatsRepositoryTests(t *testing.T, newRepository func() domain.ChatsRepository) {
	t.Run("List", func(t *testing.T) {
		t.Run("из пустого репозитория вернется пустой список", func(t *testing.T) {
			r := newRepository()
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			assert.Len(t, chats, 0)
		})
		t.Run("без фильтра из репозитория вернутся все сохраненные чаты", func(t *testing.T) {
			r := newRepository()
			chats := make([]domain.Chat, 10)
			for i := range chats {
				chats[i] = domain.Chat{
					ID:          uuid.NewString(),
					Name:        fmt.Sprintf("name%d", i),
					ChiefUserID: uuid.NewString(),
				}
				err := r.Save(chats[i])
				assert.NoError(t, err)
			}
			chatsFromRepo, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			assert.Len(t, chatsFromRepo, len(chats))
		})
		t.Run("после удаления в репозитории будет отсутствовать этот чат", func(t *testing.T) {
			r := newRepository()
			chat := domain.Chat{ID: uuid.NewString()}
			err := r.Save(chat)
			assert.NoError(t, err)
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			assert.Len(t, chats, 1)
			err = r.Delete(chat.ID)
			assert.NoError(t, err)
			chats, err = r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			assert.Len(t, chats, 0)
		})
		t.Run("с фильтром по id вернется сохраненный чат", func(t *testing.T) {
			r := newRepository()
			for range 10 {
				err := r.Save(domain.Chat{ID: uuid.NewString()})
				assert.NoError(t, err)
			}
			expectedChat := domain.Chat{ID: uuid.NewString()}
			assert.NoError(t, errors.Join(
				r.Save(expectedChat),
				r.Save(domain.Chat{ID: uuid.NewString()}),
				r.Save(domain.Chat{ID: uuid.NewString()}),
			))
			chats, err := r.List(domain.ChatsFilter{IDs: []string{expectedChat.ID}})
			assert.NoError(t, err)
			if assert.Len(t, chats, 1) {
				assertEqualChats(t, expectedChat, chats[0])
			}
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("нельзя сохранять чат без id", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Chat{
				ID:          "",
				Name:        "name",
				ChiefUserID: uuid.NewString(),
			})
			assert.Error(t, err)
		})
		t.Run("можно сохранять чат без без name", func(t *testing.T) {
			r := newRepository()
			chat := domain.Chat{
				ID:          uuid.NewString(),
				Name:        "",
				ChiefUserID: uuid.NewString(),
			}
			err := r.Save(chat)
			assert.NoError(t, err)
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			if assert.Len(t, chats, 1) {
				assertEqualChats(t, chat, chats[0])
			}
		})
		t.Run("можно сохранять чат без ChiefUserID", func(t *testing.T) {
			r := newRepository()
			chat := domain.Chat{
				ID:          uuid.NewString(),
				Name:        "name",
				ChiefUserID: "",
			}
			err := r.Save(chat)
			assert.NoError(t, err)
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			if assert.Len(t, chats, 1) {
				assertEqualChats(t, chat, chats[0])
			}
		})
		t.Run("сохраненный чат полностью соответствует сохраняемому", func(t *testing.T) {
			r := newRepository()
			chat := domain.Chat{
				ID:          uuid.NewString(),
				Name:        "name",
				ChiefUserID: uuid.NewString(),
			}
			err := r.Save(chat)
			assert.NoError(t, err)
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			if assert.Len(t, chats, 1) {
				assertEqualChats(t, chat, chats[0])
			}
		})
		t.Run("перезапись с новыми значениями по ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			const count = 33
			for i := range [count]int{} {
				err := r.Save(domain.Chat{
					ID:          id,
					Name:        fmt.Sprintf("name%d", i),
					ChiefUserID: uuid.NewString(),
				})
				assert.NoError(t, err)
			}
			expectedChat := domain.Chat{
				ID:          id,
				Name:        fmt.Sprintf("name%d", count+1),
				ChiefUserID: uuid.NewString(),
			}
			err := r.Save(expectedChat)
			assert.NoError(t, err)
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			if assert.Len(t, chats, 1) {
				assertEqualChats(t, expectedChat, chats[0])
			}
		})
		t.Run("name может дублироваться", func(t *testing.T) {
			r := newRepository()
			name := "name"
			const count = 2
			for range count {
				err := r.Save(domain.Chat{
					ID:   uuid.NewString(),
					Name: name,
				})
				assert.NoError(t, err)
			}
			chatsFromRepo, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			assert.Len(t, chatsFromRepo, count)
		})
	})
	t.Run("Delete", func(t *testing.T) {
		t.Run("параметр id обязательный", func(t *testing.T) {
			r := newRepository()
			err := r.Delete("")
			assert.Error(t, err)
		})
		t.Run("несуществующий id не вернет ошибку", func(t *testing.T) {
			r := newRepository()
			err := r.Delete(uuid.NewString())
			assert.NoError(t, err)
		})
		t.Run("без ошибок", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Chat{ID: id})
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			assert.Len(t, chats, 0)
		})
		t.Run("можно повторно удалять по ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Chat{ID: id})
			assert.NoError(t, err)
			for range 343 {
				err = r.Delete(id)
				assert.NoError(t, err)
			}
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			assert.Len(t, chats, 0)
		})
	})
}
