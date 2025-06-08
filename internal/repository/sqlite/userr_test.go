package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
	userrRepoTests "github.com/saime-0/nice-pea-chat/internal/domain/userr/repository_tests"
)

func TestUserrRepository(t *testing.T) {
	userrRepoTests.TestRepository(t, func() userr.Repository {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.Nil(t, err)
		assert.NotNil(t, repositoryFactory)
		repo := repositoryFactory.NewUserrRepository()
		require.NotNil(t, repo)
		return repo
	})
}
