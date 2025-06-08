package repository_tests

import (
	"testing"

	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
)

// TestRepository реализацию репозитория
func TestRepository(t *testing.T, newRepository func() sessionn.Repository) {
	t.Run("List", func(t *testing.T) {
		panic("implement me")
	})
	t.Run("Upsert", func(t *testing.T) {
		panic("implement me")
	})
}
