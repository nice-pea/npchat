package memory

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("создание объекта с дефолтными значениями", func(t *testing.T) {
		config := Config{}
		sqlim, err := Init(config)
		assert.NoError(t, err)
		assert.NotNil(t, sqlim)
	})
}
