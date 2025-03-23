package repository_tests

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func ChatsRepositoryTests(t *testing.T, newRepository func() domain.ChatsRepository) {
	t.Helper()
	t.Run("List", func(t *testing.T) {
		type testCase struct {
			name    string
			filter  domain.ChatsFilter
			wantRes []domain.Chat
			wantErr bool
		}
		testCaseConstructors := []func(*testing.T, domain.ChatsRepository) testCase{
			func(*testing.T, domain.ChatsRepository) testCase {
				return testCase{
					name:    "без фильтра в пустом репозитории",
					filter:  domain.ChatsFilter{},
					wantRes: []domain.Chat{},
					wantErr: false,
				}
			},
			func(t *testing.T, repo domain.ChatsRepository) testCase {
				chat := domain.Chat{
					ID:   uuid.NewString(),
					Name: "name",
				}
				assert.NoError(t, repo.Save(chat))
				return testCase{
					name:    "без фильтра, только один чат",
					filter:  domain.ChatsFilter{},
					wantRes: []domain.Chat{chat},
					wantErr: false,
				}
			},
		}

		for _, newTestCase := range testCaseConstructors {
			repo := newRepository()
			tc := newTestCase(t, repo)
			t.Run(tc.name, func(t *testing.T) {
				chats, err := repo.List(tc.filter)
				if (err != nil) != tc.wantErr {
					t.Errorf("List() error = %v, wantErr %v", err, tc.wantErr)
				}
				if !reflect.DeepEqual(chats, tc.wantRes) {
					t.Errorf("List() chats = %v, want %v", chats, tc.wantRes)
				}
			})
		}
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
}
