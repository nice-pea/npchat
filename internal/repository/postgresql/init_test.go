package postgresql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("создание соединения с дефолтными значениями", func(t *testing.T) {
		config := Config{}
		psql, err := InitRepositoryFactory(config)
		assert.Error(t, err)
		assert.Zero(t, psql)
	})
}
