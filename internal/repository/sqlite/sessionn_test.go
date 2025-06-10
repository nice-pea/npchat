package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
	sessionnRepoTests "github.com/saime-0/nice-pea-chat/internal/domain/sessionn/repository_tests"
)

func TestSessionnRepository(t *testing.T) {
	sessionnRepoTests.TestRepository(t, func() sessionn.Repository {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.Nil(t, err)
		assert.NotNil(t, repositoryFactory)
		repo := repositoryFactory.NewSessionnRepository()
		require.NotNil(t, repo)
		return repo
	})
}
