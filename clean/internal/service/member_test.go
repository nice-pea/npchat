package service

//
//import (
//	"testing"
//
//	"github.com/google/uuid"
//	"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/mock"
//
//	"github.com/saime-0/nice-pea-chat/internal/domain"
//	"github.com/saime-0/nice-pea-chat/internal/domain/mocks"
//)
//
//func TestMembers_Delete(t *testing.T) {
//	t.Run("успешное удаление из мокового репозитория", func(t *testing.T) {
//		id := uuid.NewString()
//		// Настройка мокового репозитория
//		repo := new(domain_mock.MembersRepositoryMock)
//		repo.On("List", domain.MembersFilter{ID: id}).
//			Return(make([]domain.Member, 1), nil)
//		repo.On("Delete", mock.Anything).Return(nil)
//		// Инициализация сервиса
//		service := &Members{
//			MembersRepo: repo,
//			History:     HistoryDummy{},
//		}
//		// Тест функции
//		assert.NoError(t, service.Delete(id))
//	})
//}
//
//func TestMembers_List(t *testing.T) {
//	t.Run("успешно найден участник в моковом репозитории", func(t *testing.T) {
//		findMember := domain.Member{
//			ID:     uuid.NewString(),
//			UserID: uuid.NewString(),
//			ChatID: uuid.NewString(),
//		}
//		// Настройка мокового репозитория
//		repo := new(domain_mock.MembersRepositoryMock)
//		repo.On("List", domain.MembersFilter{ID: findMember.ID}).
//			Return([]domain.Member{findMember}, nil)
//		// Инициализация сервиса
//		service := &Members{
//			MembersRepo: repo,
//			History:     HistoryDummy{},
//		}
//		// Тест функции
//		members, err := service.List(MembersListFilter{ID: findMember.ID})
//		assert.NoError(t, err)
//		assert.Equal(t, len(members), 1)
//		assert.Equal(t, findMember, members[0])
//	})
//}
