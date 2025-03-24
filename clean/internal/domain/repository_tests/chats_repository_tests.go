package repository_tests

import (
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
		t.Run("без фильтра, только один чат", func(t *testing.T) {
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
		t.Run("отсутствует после удаления", func(t *testing.T) {
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
	})
	t.Run("Save", func(t *testing.T) {
		type testCase struct {
			name    string
			newChat domain.Chat
			wantErr bool
		}
		testCaseConstructors := []func(*testing.T, domain.ChatsRepository) testCase{
			func(t *testing.T, repository domain.ChatsRepository) testCase {
				return testCase{
					name:    "чат без id",
					newChat: domain.Chat{ID: "", Name: "name"},
					wantErr: true,
				}
			},
			func(t *testing.T, repository domain.ChatsRepository) testCase {
				return testCase{
					name:    "чат без name",
					newChat: domain.Chat{ID: uuid.NewString(), Name: ""},
					wantErr: false,
				}
			},
			func(t *testing.T, repo domain.ChatsRepository) testCase {
				return testCase{
					name: "чат без ошибок",
					newChat: domain.Chat{
						ID:   uuid.NewString(),
						Name: "name",
					},
					wantErr: false,
				}
			},
			func(t *testing.T, repository domain.ChatsRepository) testCase {
				id := uuid.NewString()
				assert.NoError(t, repository.Save(domain.Chat{
					ID:   id,
					Name: "name",
				}))
				return testCase{
					name: "дубль ID",
					newChat: domain.Chat{
						ID:   id,
						Name: "name2",
					},
					wantErr: true,
				}
			},
			func(t *testing.T, repository domain.ChatsRepository) testCase {
				name := "name"
				assert.NoError(t, repository.Save(domain.Chat{
					ID:   uuid.NewString(),
					Name: name,
				}))
				return testCase{
					name: "дубль name",
					newChat: domain.Chat{
						ID:   uuid.NewString(),
						Name: "name",
					},
					wantErr: false,
				}
			},
		}

		for _, newTestCase := range testCaseConstructors {
			repo := newRepository()
			tc := newTestCase(t, repo)
			t.Run(tc.name, func(t *testing.T) {
				err := repo.Save(tc.newChat)
				if (err != nil) != tc.wantErr {
					t.Errorf("Save() error = %v, wantErr %v", err, tc.wantErr)
				}
			})
		}
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
