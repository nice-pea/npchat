package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite/memory"
)

func TestChats_ChatsWhereUserIsMember(t *testing.T) {
	newChatsService := func() *Chats {
		sqLiteInMemory, err := memory.Init(memory.Config{MigrationsDir: "../../migrations/repository/sqlite/memory"})
		assert.NoError(t, err)
		chatsRepository, err := sqLiteInMemory.NewChatsRepository()
		assert.NoError(t, err)
		membersRepository, err := sqLiteInMemory.NewMembersRepository()
		assert.NoError(t, err)
		return &Chats{
			ChatsRepo:   chatsRepository,
			MembersRepo: membersRepository,
			History:     HistoryDummy{},
		}
	}
	t.Run("SubjectUserID обязательное поле", func(t *testing.T) {
		input := ChatsWhereUserIsMemberInput{
			SubjectUserID: "",
			UserID:        uuid.NewString(),
		}
		userChats, err := newChatsService().ChatsWhereUserIsMember(input)
		assert.Error(t, err)
		assert.Len(t, userChats, 0)
	})
	t.Run("UserID обязательное поле", func(t *testing.T) {
		input := ChatsWhereUserIsMemberInput{
			SubjectUserID: uuid.NewString(),
			UserID:        "",
		}
		userChats, err := newChatsService().ChatsWhereUserIsMember(input)
		assert.Error(t, err)
		assert.Len(t, userChats, 0)
	})
	t.Run("UserID и SubjectUserID разные", func(t *testing.T) {
		input := ChatsWhereUserIsMemberInput{
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		userChats, err := newChatsService().ChatsWhereUserIsMember(input)
		assert.Error(t, err)
		assert.Len(t, userChats, 0)
	})
	t.Run("пустой список из пустого репозитория", func(t *testing.T) {
		id := uuid.NewString()
		input := ChatsWhereUserIsMemberInput{
			SubjectUserID: id,
			UserID:        id,
		}
		userChats, err := newChatsService().ChatsWhereUserIsMember(input)
		assert.NoError(t, err)
		assert.Len(t, userChats, 0)
	})
	t.Run("пустой список из заполненного репозитория", func(t *testing.T) {
		chatsService := newChatsService()
		for i := 0; i < 7; i++ {
			// Создать чат
			chatID := uuid.NewString()
			err := chatsService.ChatsRepo.Save(domain.Chat{
				ID:          chatID,
				Name:        "name" + string(rune(i)),
				ChiefUserID: uuid.NewString(),
			})
			assert.NoError(t, err)
			// Создать участника в чате
			err = chatsService.MembersRepo.Save(domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chatID,
			})
			assert.NoError(t, err)
		}
		userID := uuid.NewString()
		input := ChatsWhereUserIsMemberInput{
			SubjectUserID: userID,
			UserID:        userID,
		}
		userChats, err := chatsService.ChatsWhereUserIsMember(input)
		assert.NoError(t, err)
		assert.Len(t, userChats, 0)
	})
	t.Run("3 чата из заполненного репозитория", func(t *testing.T) {
		chatsService := newChatsService()
		var existsChats []domain.Chat
		for i := 0; i < 10; i++ {
			// Создать чат
			existsChats = append(existsChats, domain.Chat{
				ID:          uuid.NewString(),
				Name:        "name" + string(rune(i)),
				ChiefUserID: uuid.NewString(),
			})
			err := chatsService.ChatsRepo.Save(existsChats[i])
			assert.NoError(t, err)
			// Создать участника в чате
			err = chatsService.MembersRepo.Save(domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: existsChats[i].ID,
			})
			assert.NoError(t, err)
		}
		userID := uuid.NewString()
		for i := 0; i < 3; i++ {
			err := chatsService.MembersRepo.Save(domain.Member{
				ID:     uuid.NewString(),
				UserID: userID,
				ChatID: existsChats[i].ID,
			})
			assert.NoError(t, err)
		}
		input := ChatsWhereUserIsMemberInput{
			SubjectUserID: userID,
			UserID:        userID,
		}
		userChats, err := chatsService.ChatsWhereUserIsMember(input)
		assert.NoError(t, err)
		assert.Len(t, userChats, 3)
	})
}

func Test_CreateChat(t *testing.T) {
	//assertChatEqualIn := func(in CreateChatInput, out domain.Chat) {
	//	assert.Equal(t, out.ChiefUserID, in.ChiefUserID)
	//	assert.Equal(t, out.Name, in.Name)
	//}
	//t.Run("удачное создание чата", func(t *testing.T) {
	//	s := Chats{}
	//	input := CreateChatInput{
	//		ChiefUserID: uuid.NewString(),
	//		Name:        "abcd",
	//	}
	//	chat, err := s.Create(input)
	//	assert.NoError(t, err)
	//	assert.NotZero(t, chat)
	//})
	//t.Run("выходящие совпадают с заданными", func(t *testing.T) {
	//	s := Chats{}
	//	input := CreateChatInput{
	//		ChiefUserID: uuid.NewString(),
	//		Name:        "wqert",
	//	}
	//	chat, err := s.Create(input)
	//	assert.NoError(t, err)
	//	assertChatEqualIn(input, chat)
	//})
	//t.Run("ошибка при пустом name", func(t *testing.T) {
	//	s := Chats{}
	//	input := CreateChatInput{
	//		ChiefUserID: uuid.NewString(),
	//		Name:        "",
	//	}
	//	chat, err := s.Create(input)
	//	assert.Error(t, err)
	//	assert.Zero(t, chat)
	//})
	//t.Run("ошибка при пустом ChiefUserID", func(t *testing.T) {
	//	s := Chats{}
	//	input := CreateChatInput{
	//		ChiefUserID: "",
	//		Name:        "qf",
	//	}
	//	chat, err := s.Create(input)
	//	assert.Error(t, err)
	//	assert.Zero(t, chat)
	//})
	//t.Run("создать 2 раза с одинаковыми параметрами", func(t *testing.T) {
	//	s := Chats{}
	//	input := CreateChatInput{ChiefUserID: uuid.NewString(), Name: "sda"}
	//	for range [2]int{} {
	//		chat, err := s.Create(input)
	//		assert.NoError(t, err)
	//		assertChatEqualIn(input, chat)
	//	}
	//})
	//t.Run("достичь лимита", func(t *testing.T) {
	//	s := Chats{}
	//	input := CreateChatInput{ChiefUserID: uuid.NewString(), Name: "asdgge"}
	//	for range [LimitCreatedChats]int{} {
	//		chat, err := s.Create(input)
	//		assert.NoError(t, err)
	//		assertChatEqualIn(input, chat)
	//	}
	//	chat, err := s.Create(input)
	//	assert.Error(t, err)
	//	assert.Zero(t, chat)
	//})
}
