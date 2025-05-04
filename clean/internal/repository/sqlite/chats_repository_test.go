package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/repository_tests"
)

func TestNewChatsRepository(t *testing.T) {
	t.Run("создание репозитория с дефолтными значениями", func(t *testing.T) {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.Nil(t, err)
		assert.NotNil(t, repositoryFactory)
		repo := repositoryFactory.NewChatsRepository()
		assert.NotNil(t, repo)
	})
}

func TestChatsRepository(t *testing.T) {
	repository_tests.ChatsRepositoryTests(t, func() domain.ChatsRepository {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.Nil(t, err)
		assert.NotNil(t, repositoryFactory)
		repo := repositoryFactory.NewChatsRepository()
		require.NotNil(t, repo)
		return repo
	})
}
