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
		t.Run("без фильтра в пустом репозитории", func(t *testing.T) {
			r := newRepository()
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			assert.Len(t, chats, 0)
		})
		t.Run("без фильтра только один чат", func(t *testing.T) {
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
		t.Run("без фильтра отсутствует после удаления", func(t *testing.T) {
			r := newRepository()
			chat := domain.Chat{
				ID:   uuid.NewString(),
				Name: "name",
			}
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
		t.Run("с фильтром по id", func(t *testing.T) {
			r := newRepository()
			chat := domain.Chat{ID: uuid.NewString()}
			assert.NoError(t, errors.Join(
				r.Save(chat),
				r.Save(domain.Chat{ID: uuid.NewString()}),
				r.Save(domain.Chat{ID: uuid.NewString()}),
			))
			chats, err := r.List(domain.ChatsFilter{IDs: []string{chat.ID}})
			assert.NoError(t, err)
			if assert.Len(t, chats, 1) {
				assertEqualChats(t, chat, chats[0])
			}
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("чат без id", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Chat{
				ID:          "",
				Name:        "name",
				ChiefUserID: uuid.NewString(),
			})
			assert.Error(t, err)
		})
		t.Run("чат без name", func(t *testing.T) {
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
		t.Run("чат без ChiefUserID", func(t *testing.T) {
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
		t.Run("с полными данными", func(t *testing.T) {
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
		t.Run("дубль ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Chat{
				ID:          id,
				Name:        "name",
				ChiefUserID: uuid.NewString(),
			})
			assert.NoError(t, err)
			err = r.Save(domain.Chat{
				ID:          id,
				Name:        "name1",
				ChiefUserID: uuid.NewString(),
			})
			assert.Error(t, err)
		})
		t.Run("дубль name", func(t *testing.T) {
			r := newRepository()
			const count = 2
			chats := make([]domain.Chat, count)
			for i := range [count]int{} {
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
			assert.Len(t, chatsFromRepo, count)
			for i, chat := range chats {
				assertEqualChats(t, chat, chatsFromRepo[i])
			}
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
			err := r.Save(domain.Chat{
				ID:   id,
				Name: "name",
			})
			err = r.Delete(id)
			assert.NoError(t, err)
		})
		t.Run("дважды удаленный", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Chat{
				ID:   id,
				Name: "name",
			})
			err = r.Delete(id)
			assert.NoError(t, err)
			err = r.Delete(id)
			assert.NoError(t, err)
		})
	})
}
