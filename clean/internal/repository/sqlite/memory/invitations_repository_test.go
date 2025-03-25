package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/repository_tests"
)

func TestNewInvitationsRepository(t *testing.T) {
	t.Run("создание репозитория с дефолтными значениями", func(t *testing.T) {
		sqlim, err := Init(Config{
			MigrationsDir: "../../../../migrations/repository/sqlite/memory",
		})
		assert.Nil(t, err)
		assert.NotNil(t, sqlim)
		repo, err := sqlim.NewInvitationsRepository()
		assert.Nil(t, err)
		assert.NotNil(t, repo)
	})
}

func TestInvitationsRepository(t *testing.T) {
	repository_tests.InvitationsRepositoryTests(t, func() domain.InvitationsRepository {
		sqlim, err := Init(Config{
			MigrationsDir: "../../../../migrations/repository/sqlite/memory",
		})
		assert.Nil(t, err)
		assert.NotNil(t, sqlim)
		repo, err := sqlim.NewInvitationsRepository()
		assert.Nil(t, err)
		assert.NotNil(t, repo)
		return repo
	})
}
