package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/repository_tests"
)

func TestNewSessionsRepository(t *testing.T) {
	t.Run("создание репозитория с дефолтными значениями", func(t *testing.T) {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.NoError(t, err)
		require.NotZero(t, repositoryFactory)
		repo := repositoryFactory.NewSessionsRepository()
		assert.NotNil(t, repo)
	})
}

func TestSessionsRepository(t *testing.T) {
	repository_tests.SessionsRepositoryTests(t, func() domain.SessionsRepository {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.NoError(t, err)
		require.NotZero(t, repositoryFactory)
		repo := repositoryFactory.NewSessionsRepository()
		require.NotNil(t, repo)
		return repo
	})
}
