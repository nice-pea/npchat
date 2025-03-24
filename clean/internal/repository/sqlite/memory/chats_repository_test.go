package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/repository_tests"
)

// todo:
func TestNewChatsRepository(t *testing.T) {
	t.Run("создание репозитория с дефолтными значениями", func(t *testing.T) {
		sqlim, err := Init(Config{
			MigrationsDir: "../../../../migrations/repository/sqlite/memory",
		})
		assert.Nil(t, err)
		assert.NotNil(t, sqlim)
		repo, err := sqlim.NewChatsRepository()
		assert.Nil(t, err)
		assert.NotNil(t, repo)
	})
}

func TestChatsRepository(t *testing.T) {
	repository_tests.ChatsRepositoryTests(t, func() domain.ChatsRepository {
		sqlim, err := Init(Config{
			MigrationsDir: "../../../../migrations/repository/sqlite/memory",
		})
		assert.Nil(t, err)
		assert.NotNil(t, sqlim)
		repo, err := sqlim.NewChatsRepository()
		assert.Nil(t, err)
		assert.NotNil(t, repo)
		return repo
	})
}
