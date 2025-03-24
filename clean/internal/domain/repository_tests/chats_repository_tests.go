package repository_tests

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func ChatsRepositoryTests(t *testing.T, newRepository func() domain.ChatsRepository) {
	t.Helper()
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
				ID:   uuid.NewString(),
				Name: "name",
			}
			err := r.Save(chat)
			assert.NoError(t, err)
			chats, err := r.List(domain.ChatsFilter{})
			assert.NoError(t, err)
			assert.Len(t, chats, 1)
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
			id := uuid.NewString()
			assert.NoError(t, errors.Join(
				r.Save(domain.Chat{ID: id, Name: "name"}),
				r.Save(domain.Chat{ID: uuid.NewString(), Name: "name1"}),
				r.Save(domain.Chat{ID: uuid.NewString(), Name: "name2"}),
			))
			chats, err := r.List(domain.ChatsFilter{ID: id})
			assert.NoError(t, err)
			assert.Len(t, chats, 1)
		})
	})
	t.Run("Save", func(t *testing.T) {
		t.Run("чат без id", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Chat{
				ID:   "",
				Name: "name",
			})
			assert.Error(t, err)
		})
		t.Run("чат без name", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Chat{
				ID:   uuid.NewString(),
				Name: "",
			})
			assert.NoError(t, err)
		})
		t.Run("без ошибок", func(t *testing.T) {
			r := newRepository()
			err := r.Save(domain.Chat{
				ID:   uuid.NewString(),
				Name: "name",
			})
			assert.NoError(t, err)
		})
		t.Run("дубль ID", func(t *testing.T) {
			r := newRepository()
			id := uuid.NewString()
			err := r.Save(domain.Chat{
				ID:   id,
				Name: "name",
			})
			assert.NoError(t, err)
			err = r.Save(domain.Chat{
				ID:   id,
				Name: "name1",
			})
			assert.Error(t, err)
		})
		t.Run("дубль name", func(t *testing.T) {
			r := newRepository()
			name := "name"
			err := r.Save(domain.Chat{
				ID:   uuid.NewString(),
				Name: name,
			})
			assert.NoError(t, err)
			err = r.Save(domain.Chat{
				ID:   uuid.NewString(),
				Name: name,
			})
			assert.NoError(t, err)
		})
	})
	t.Run("Delete", func(t *testing.T) {
		type testCase struct {
			name    string
			chatID  string
			wantErr bool
		}
		testCaseConstructors := []func(*testing.T, domain.ChatsRepository) testCase{
			func(t *testing.T, repository domain.ChatsRepository) testCase {
				return testCase{
					name:    "с пустым id",
					chatID:  "",
					wantErr: true,
				}
			},
			func(t *testing.T, repository domain.ChatsRepository) testCase {
				return testCase{
					name:    "несуществующий id",
					chatID:  "c2e93bd8-dc78-4e9c-876e-07130c0b0224",
					wantErr: false,
				}
			},
			func(t *testing.T, repository domain.ChatsRepository) testCase {
				id := "c2e93bd8-dc78-4e9c-876e-07130c0b0224"
				assert.NoError(t, repository.Delete(id))
				return testCase{
					name:    "дважды удаленный",
					chatID:  "c2e93bd8-dc78-4e9c-876e-07130c0b0224",
					wantErr: false,
				}
			},
		}

		for _, newTestCase := range testCaseConstructors {
			repo := newRepository()
			tc := newTestCase(t, repo)
			t.Run(tc.name, func(t *testing.T) {
				err := repo.Delete(tc.chatID)
				if (err != nil) != tc.wantErr {
					t.Errorf("Delete() error = %v, wantErr %v", err, tc.wantErr)
				}
			})
		}
	})
}
