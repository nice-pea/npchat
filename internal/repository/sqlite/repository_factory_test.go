package sqlite

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInitRepositoryFactory(t *testing.T) {
	t.Run("создание объекта с дефолтными значениями", func(t *testing.T) {
		config := Config{}
		sqlim, err := InitRepositoryFactory(config)
		assert.NoError(t, err)
		assert.NotNil(t, sqlim)
	})
}

func TestRepositoryFactoryNewRepository(t *testing.T) {
	t.Run("создание репозиториев с дефолтными значениями", func(t *testing.T) {
		repositoryFactory, err := InitRepositoryFactory(defaultTestConfig)
		assert.NoError(t, err)
		require.NotZero(t, repositoryFactory)
		assert.NotZero(t, repositoryFactory.NewUserrRepository())
		assert.NotZero(t, repositoryFactory.NewSessionnRepository())
		assert.NotZero(t, repositoryFactory.NewChattRepository())
	})
}
