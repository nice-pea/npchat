package memory

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/repository_tests"
)

func TestNewUsersRepository(t *testing.T) {
	t.Run("создание репозитория с дефолтными значениями", func(t *testing.T) {
		sqlim, err := Init(Config{
			MigrationsDir: "../../../../migrations/repository/sqlite/memory",
		})
		assert.Nil(t, err)
		assert.NotNil(t, sqlim)
		repo, err := sqlim.NewUsersRepository()
		assert.Nil(t, err)
		assert.NotNil(t, repo)
	})
}

func TestUsersRepository(t *testing.T) {
	repository_tests.UsersRepositoryTests(t, func() domain.UsersRepository {
		sqlim, err := Init(Config{
			MigrationsDir: "../../../../migrations/repository/sqlite/memory",
		})
		assert.Nil(t, err)
		assert.NotNil(t, sqlim)
		repo, err := sqlim.NewUsersRepository()
		assert.Nil(t, err)
		assert.NotNil(t, repo)
		return repo
	})
}

func TestUsersRepository_Mapping(t *testing.T) {
	t.Run("один в domain", func(t *testing.T) {
		repoUser := user{
			ID: uuid.NewString(),
		}
		domainUser := userToDomain(repoUser)
		assert.Equal(t, repoUser.ID, domainUser.ID)
	})
	t.Run("один из domain", func(t *testing.T) {
		domainUser := domain.User{
			ID: uuid.NewString(),
		}
		repoUser := userFromDomain(domainUser)
		assert.Equal(t, domainUser.ID, repoUser.ID)
	})
	t.Run("несколько в domain", func(t *testing.T) {
		repoUsers := []user{
			{ID: uuid.NewString()},
			{ID: uuid.NewString()},
			{ID: uuid.NewString()},
		}
		domainUsers := usersToDomain(repoUsers)
		for i, repoUser := range repoUsers {
			assert.Equal(t, repoUser.ID, domainUsers[i].ID)
		}
	})
	t.Run("несколько из domain", func(t *testing.T) {
		domainUsers := []domain.User{
			{ID: uuid.NewString()},
			{ID: uuid.NewString()},
			{ID: uuid.NewString()},
		}
		repoUsers := usersFromDomain(domainUsers)
		for i, domainUser := range domainUsers {
			assert.Equal(t, domainUser.ID, repoUsers[i].ID)
		}
	})
}
