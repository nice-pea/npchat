package base

import (
	"testing"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/repository_tests"
)

func TestChatsRepository(t *testing.T) {
	constructor := func(t *testing.T) domain.ChatsRepository {
		// todo: написать инициализацию базового репозитория, но в тестовом окружении
		return &ChatsRepository{}
	}
	//repository_tests.TestChatsRepositoryDelete(t, constructor)
	repository_tests.ChatsRepositoryList(t, constructor)
	//repository_tests.TestChatsRepositorySave(t, constructor)
}
