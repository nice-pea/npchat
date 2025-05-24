package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/repository_tests"
)

func TestNewOAuthRepository(t *testing.T) {
	t.Run("создание репозитория с дефолтными значениями", func(t *testing.T) {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.NoError(t, err)
		require.NotZero(t, repositoryFactory)
		repo := repositoryFactory.NewOAuthRepository()
		assert.NotNil(t, repo)
	})
}

func TestOAuthRepository(t *testing.T) {
	repository_tests.OAuthRepositoryTests(t, func() domain.OAuthRepository {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.NoError(t, err)
		require.NotZero(t, repositoryFactory)
		repo := repositoryFactory.NewOAuthRepository()
		require.NotNil(t, repo)
		return repo
	})
}
