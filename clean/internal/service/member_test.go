package service

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/mocks"
)

func TestMembers_Delete(t *testing.T) {
	t.Run("ok", func(t *testing.T) {
		id := uint(rand.Int())
		// Настройка мокового репозитория
		repo := new(domain_mock.MembersRepositoryMock)
		repo.On("List", domain.MembersFilter{ID: id}).
			Return(make([]domain.Member, 1), nil)
		repo.On("Delete", mock.Anything).Return(nil)
		// Инициализация сервиса
		service := &Members{
			MembersRepo: repo,
			Hi:          HistoryDummy{},
		}
		// Тест функции
		assert.NoError(t, service.Delete(id), repo)
	})
}
