package service

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite/memory"
)

// newChatsService создает объект сервиса Chats с sqlite/memory репозиториями
func newMembersService(t *testing.T) *Chats {
	sqLiteInMemory, err := memory.Init(memory.Config{MigrationsDir: "../../migrations/repository/sqlite/memory"})
	assert.NoError(t, err)
	chatsRepository, err := sqLiteInMemory.NewChatsRepository()
	assert.NoError(t, err)
	membersRepository, err := sqLiteInMemory.NewMembersRepository()
	assert.NoError(t, err)
	return &Members{
		//ChatsRepo:   chatsRepository,
		//MembersRepo: membersRepository,
		//History:     HistoryDummy{},
	}
}

/*
Создать сервис service/Members и добавить ему несколько юзкейсов:
Получить список участников чата

	доступно только участникам чата
	в список попадают все активные участники чата
	Входящие параметры должны валидироваться

Покинуть чат

	недоступно для chief
	доступно только участникам чата
	Входящие параметры должны валидироваться

Принудительно удалить участника

	доступно только для chief
	нельзя удалить самого себя
	Входящие параметры должны валидироваться
*/

func Test_Members_ChatMembers(t *testing.T) {

}

func Test_Members_Leave(t *testing.T) {

}
func Test_Members_Delete(t *testing.T) {

}
