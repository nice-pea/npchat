package sessionn

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

// НАСТРОЙКА ТЕСТОВОГО ОКРУЖЕНИЯ

type repoSuite struct {
	suite.Suite
	newRepository func() Repository
	repo          Repository
}

// Test_repoSuite запускает набор тестов
func Test_repoSuite(t *testing.T) {
	suite.Run(t, &repoSuite{
		Suite:         suite.Suite{},
		newRepository: sqliteRepoConstructor(t),
		repo:          nil,
	})
}

// SetupSubTest подготавливает репозиторий для каждого подтеста
func (suite *repoSuite) SetupSubTest() {
	suite.repo = suite.newRepository()
}

func (suite *repoSuite) TearDownSubTest() {}

// sqliteRepoConstructor создает конструктор SQLite репозитория с тестовой конфигурацией
func sqliteRepoConstructor(t *testing.T) func() Repository {
	sqliteConfig := sqlite.Config{
		MigrationsDir: "../../../migrations/repository/sqlite",
	}
	repositoryFactory, err := sqlite.InitRepositoryFactory(sqliteConfig)
	assert.Nil(t, err)
	assert.NotNil(t, repositoryFactory)
	repo := repositoryFactory.NewChatsRepository()
	require.NotNil(t, repo)

	f, err := sqlite.InitRepositoryFactory(sqliteConfig)
	require.NoError(t, err)
	return f.NewSessionnRepository
}

// ТЕСТЫ

// Test_Repository реализацию репозитория
func (suite *repoSuite) Test_Repository() {
	suite.Run("List", func() {
		panic("implement me")
	})
	suite.Run("Upsert", func() {
		panic("implement me")
	})

}
