package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/repository_tests"
)

func TestNewLoginCredentialsRepository(t *testing.T) {
	t.Run("создание репозитория с дефолтными значениями", func(t *testing.T) {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.NoError(t, err)
		require.NotZero(t, repositoryFactory)
		repo := repositoryFactory.NewLoginCredentialsRepository()
		assert.NotZero(t, repo)
	})
}

func TestLoginCredentialsRepository(t *testing.T) {
	repository_tests.LoginCredentialsRepositoryTests(t, func() domain.LoginCredentialsRepository {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.NoError(t, err)
		require.NotZero(t, repositoryFactory)
		repo := repositoryFactory.NewLoginCredentialsRepository()
		require.NotZero(t, repo)
		return repo
	})
}
