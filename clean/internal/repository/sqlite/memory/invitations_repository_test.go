package memory

import (
	"testing"

	"github.com/google/uuid"
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

func TestInvitationsRepository_Mapping(t *testing.T) {
	t.Run("один в domain", func(t *testing.T) {
		repoInvitation := invitation{
			ID:     uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		domainInvitations := invitationToDomain(repoInvitation)
		assert.Equal(t, repoInvitation.ID, domainInvitations.ID)
		assert.Equal(t, repoInvitation.ChatID, domainInvitations.ChatID)
	})
	t.Run("один из domain", func(t *testing.T) {
		domainInvitations := domain.Invitation{
			ID:     uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		repoInvitation := invitationFromDomain(domainInvitations)
		assert.Equal(t, domainInvitations.ID, repoInvitation.ID)
		assert.Equal(t, domainInvitations.ChatID, repoInvitation.ChatID)
	})
	t.Run("несколько в domain", func(t *testing.T) {
		repoInvitations := []invitation{
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
		}
		domainInvitations := invitationsToDomain(repoInvitations)
		for i, repoInvitation := range repoInvitations {
			assert.Equal(t, repoInvitation.ID, domainInvitations[i].ID)
			assert.Equal(t, repoInvitation.ChatID, domainInvitations[i].ChatID)
		}
	})
	t.Run("несколько из domain", func(t *testing.T) {
		domainInvitations := []domain.Invitation{
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
			{ID: uuid.NewString(), ChatID: uuid.NewString()},
		}
		repoInvitations := invitationsFromDomain(domainInvitations)
		for i, domainInvitation := range domainInvitations {
			assert.Equal(t, domainInvitation.ID, repoInvitations[i].ID)
			assert.Equal(t, domainInvitation.ChatID, repoInvitations[i].ChatID)
		}
	})
}
