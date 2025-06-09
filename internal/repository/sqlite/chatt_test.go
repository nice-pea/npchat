package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
	chattRepoTests "github.com/saime-0/nice-pea-chat/internal/domain/chatt/repository_tests"
)

func TestChattRepository(t *testing.T) {
	chattRepoTests.TestRepository(t, func() chatt.Repository {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.Nil(t, err)
		assert.NotNil(t, repositoryFactory)
		repo := repositoryFactory.NewChattRepository()
		require.NotNil(t, repo)
		return repo
	})
}
