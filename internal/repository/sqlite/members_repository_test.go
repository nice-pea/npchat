package sqlite

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/repository_tests"
)

func TestNewMembersRepository(t *testing.T) {
	t.Run("создание репозитория с дефолтными значениями", func(t *testing.T) {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.Nil(t, err)
		assert.NotNil(t, repositoryFactory)
		repo := repositoryFactory.NewMembersRepository()
		assert.NotNil(t, repo)
	})
}

func TestMembersRepository(t *testing.T) {
	repository_tests.MembersRepositoryTests(t, func() domain.MembersRepository {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.Nil(t, err)
		assert.NotNil(t, repositoryFactory)
		repo := repositoryFactory.NewMembersRepository()
		require.NotNil(t, repo)
		return repo
	})
}

func TestMembersRepository_Mapping(t *testing.T) {
	t.Run("один в domain", func(t *testing.T) {
		repoMember := member{
			ID:     uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		domainMember := memberToDomain(repoMember)
		assert.Equal(t, repoMember.ID, domainMember.ID)
		assert.Equal(t, repoMember.ChatID, domainMember.ChatID)
	})
	t.Run("один из domain", func(t *testing.T) {
		domainMember := domain.Member{
			ID:     uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		repoMember := memberFromDomain(domainMember)
		assert.Equal(t, domainMember.ID, repoMember.ID)
		assert.Equal(t, domainMember.ChatID, repoMember.ChatID)
	})
	t.Run("несколько в domain", func(t *testing.T) {
		repoMembers := []member{
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
		}
		domainMembers := membersToDomain(repoMembers)
		for i, repoMember := range repoMembers {
			assert.Equal(t, repoMember.ID, domainMembers[i].ID)
			assert.Equal(t, repoMember.ChatID, domainMembers[i].ChatID)
		}
	})
	t.Run("несколько из domain", func(t *testing.T) {
		domainMembers := []domain.Member{
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
		}
		repoMembers := membersFromDomain(domainMembers)
		for i, domainMember := range domainMembers {
			assert.Equal(t, domainMember.ID, repoMembers[i].ID)
			assert.Equal(t, domainMember.ChatID, repoMembers[i].ChatID)
		}
	})
}
