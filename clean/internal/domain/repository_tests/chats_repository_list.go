package repository_tests

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

func ChatsRepositoryList(t *testing.T, newRepository func(*testing.T) domain.ChatsRepository) {
	t.Helper()
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
		repo := newRepository(t)
		tc := newTestCase(t, repo)
		t.Run(tc.name, func(t *testing.T) {
			chats, err := repo.List(tc.filter)
			if (err != nil) != tc.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tc.wantErr)
			}
			if !reflect.DeepEqual(chats, tc.wantRes) {
				t.Errorf("List() chats = %v, want %v", chats, tc.filter)
			}
		})
	}
}
