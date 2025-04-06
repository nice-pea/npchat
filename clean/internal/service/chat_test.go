package service

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite/memory"
)

type chatsTestEnv struct {
	chatsService *Chats
	t            *testing.T
}

func initChatTestEnv(t *testing.T) chatsTestEnv {
	env := chatsTestEnv{
		chatsService: &Chats{},
		t:            t,
	}
	sqLiteInMemory, err := memory.Init(memory.Config{MigrationsDir: "../../migrations/repository/sqlite/memory"})
	assert.NoError(t, err)
	env.chatsService.ChatsRepo, err = sqLiteInMemory.NewChatsRepository()
	assert.NoError(t, err)
	env.chatsService.MembersRepo, err = sqLiteInMemory.NewMembersRepository()
	assert.NoError(t, err)

	return env
}

func (e *chatsTestEnv) newUserChatsInput(userID string) UserChatsInput {
	return UserChatsInput{
		SubjectUserID: userID,
		UserID:        userID,
	}
}

func (e *chatsTestEnv) saveChat(chat domain.Chat) domain.Chat {
	err := e.chatsService.ChatsRepo.Save(chat)
	assert.NoError(e.t, err)

	return chat
}

func (e *chatsTestEnv) saveMember(member domain.Member) domain.Member {
	err := e.chatsService.MembersRepo.Save(member)
	assert.NoError(e.t, err)

	return member
}

func (e *chatsTestEnv) assertChatEqualInput(in CreateInput, chat domain.Chat) {
	assert.Equal(e.t, in.Name, chat.Name)
	assert.Equal(e.t, in.ChiefUserID, chat.ChiefUserID)
}
func (e *chatsTestEnv) newCreateInputRandom() CreateInput {
	return CreateInput{
		ChiefUserID: uuid.NewString(),
		Name:        fmt.Sprintf("name%d", rand.Int()),
	}
}

// Test_UserChatsInput_Validate тестирует валидацию входящих параметров запроса списка чатов в которых участвует пользователь
func Test_UserChatsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := UserChatsInput{
			SubjectUserID: id,
			UserID:        id,
		}
		return in.Validate()
	})
}

// Test_Chats_UserChats тестирует запрос список чатов в которых участвует пользователь
func Test_Chats_UserChats(t *testing.T) {
	t.Run("пользователь может запрашивать только свой чат", func(t *testing.T) {
		env := initChatTestEnv(t)
		input := UserChatsInput{
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		chats, err := env.chatsService.UserChats(input)
		assert.ErrorIs(t, err, ErrUnauthorizedChatsView)
		assert.Len(t, chats, 0)
	})
	t.Run("пустой список из пустого репозитория", func(t *testing.T) {
		env := initChatTestEnv(t)
		input := env.newUserChatsInput(uuid.NewString())
		userChats, err := env.chatsService.UserChats(input)
		assert.NoError(t, err)
		assert.Len(t, userChats, 0)
	})
	t.Run("пустой список если у пользователя нет чатов", func(t *testing.T) {
		env := initChatTestEnv(t)
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := env.saveChat(domain.Chat{
				ID: uuid.NewString(),
			})
			env.saveMember(domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			})
		}
		input := env.newUserChatsInput(uuid.NewString())
		userChats, err := env.chatsService.UserChats(input)
		assert.NoError(t, err)
		assert.Len(t, userChats, 0)
	})
	t.Run("у пользователя может быть несколько чатов", func(t *testing.T) {
		env := initChatTestEnv(t)
		userID := uuid.NewString()
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := env.saveChat(domain.Chat{
				ID: uuid.NewString(),
			})
			env.saveMember(domain.Member{
				ID:     uuid.NewString(),
				UserID: userID,
				ChatID: chat.ID,
			})
		}
		input := env.newUserChatsInput(userID)
		userChats, err := env.chatsService.UserChats(input)
		assert.NoError(t, err)
		assert.Len(t, userChats, chatsAllCount)
	})
}

// Test_CreateChatInput_Validate тестирует валидацию входящих параметров для создания чата
func Test_CreateChatInput_Validate(t *testing.T) {
	t.Run("ошибка при пустом name", func(t *testing.T) {
		input := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "",
		}
		assert.ErrorIs(t, input.Validate(), ErrInvalidName)
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := CreateInput{
			ChiefUserID: id,
			Name:        "validName",
		}
		return in.Validate()
	})
}

// Test_Chats_CreateChat тестирует создание чата
func Test_Chats_CreateChat(t *testing.T) {
	t.Run("выходящие совпадают с заданными", func(t *testing.T) {
		env := initChatTestEnv(t)
		// Создать чат
		input := env.newCreateInputRandom()
		out, err := env.chatsService.Create(input)
		assert.NoError(t, err)
		// Сравнить результат с входящими значениями
		assert.Equal(t, input.ChiefUserID, out.ChiefMember.UserID)
		env.assertChatEqualInput(input, out.Chat)
	})
	t.Run("можно затем прочитать из репозитория", func(t *testing.T) {
		env := initChatTestEnv(t)
		// Создать чат
		input := env.newCreateInputRandom()
		out, err := env.chatsService.Create(input)
		assert.NoError(t, err)
		assert.NotZero(t, out)
		// Получить список чатов
		chats, err := env.chatsService.ChatsRepo.List(domain.ChatsFilter{})
		assert.NoError(t, err)
		// В списке этот чат будет единственным
		if assert.Len(t, chats, 1) {
			env.assertChatEqualInput(input, chats[0])
		}
	})
	t.Run("создается участник для главного администратора", func(t *testing.T) {
		env := initChatTestEnv(t)
		// Создать чат
		input := env.newCreateInputRandom()
		out, err := env.chatsService.Create(input)
		assert.NoError(t, err)
		assert.NotZero(t, out)
		// Получить список участников
		members, err := env.chatsService.MembersRepo.List(domain.MembersFilter{})
		assert.NoError(t, err)
		// В списке этот участник будет единственным
		if assert.Len(t, members, 1) {
			// Участником является главный администратор созданного чата
			assert.Equal(t, input.ChiefUserID, members[0].UserID)
			assert.Equal(t, out.Chat.ID, members[0].ChatID)
		}
	})
	t.Run("можно создать чаты с одинаковым именем", func(t *testing.T) {
		env := initChatTestEnv(t)
		input := env.newCreateInputRandom()
		// Создать несколько чатов с одинаковым именем
		const chatsAllCount = 2
		for range chatsAllCount {
			out, err := env.chatsService.Create(input)
			assert.NoError(t, err)
			assert.NotZero(t, out)
		}
		// Получить список чатов
		chats, err := env.chatsService.ChatsRepo.List(domain.ChatsFilter{})
		assert.NoError(t, err)
		// Количество чатов равно количеству созданных
		assert.Len(t, chats, chatsAllCount)
	})
	t.Run("количество созданных чатов на одного пользователя не ограничено", func(t *testing.T) {
		env := initChatTestEnv(t)
		// Пользователь
		userID := uuid.NewString()
		// Создать много чатов от лица пользователя
		const chatsAllCount = 900
		for range chatsAllCount {
			out, err := env.chatsService.Create(CreateInput{
				ChiefUserID: userID,
				Name:        "name",
			})
			assert.NoError(t, err)
			assert.NotZero(t, out)
		}
		// Получить список чатов
		chats, err := env.chatsService.ChatsRepo.List(domain.ChatsFilter{})
		assert.NoError(t, err)
		// Количество чатов равно количеству созданных
		assert.Len(t, chats, chatsAllCount)
	})
}

// Test_UpdateNameInput_Validate тестирует входящие параметры обновления названия чата
func Test_UpdateNameInput_Validate(t *testing.T) {
	t.Run("NewName", func(t *testing.T) {
		tests := []struct {
			name    string
			newName string
			wantErr bool
		}{
			{
				name:    "пустая строка",
				newName: "",
				wantErr: true,
			},
			{
				name:    "превышает лимит в 50 символов",
				newName: strings.Repeat("a", 51),
				wantErr: true,
			},
			{
				name:    "содержит пробел в начале",
				newName: " name",
				wantErr: true,
			},
			{
				name:    "содержит пробел в конце",
				newName: "name ",
				wantErr: true,
			},
			{
				name:    "содержит таб",
				newName: "na\tme",
				wantErr: true,
			},
			{
				name:    "содержит новую строку",
				newName: "na\nme",
				wantErr: true,
			},
			{
				name:    "содержит цифры",
				newName: "1na13me4",
				wantErr: false,
			},
			{
				name:    "содержит пробел в середине",
				newName: "na me",
				wantErr: false,
			},
			{
				name:    "содержит пробелы в середине",
				newName: "na  me",
				wantErr: false,
			},
			{
				name:    "содержит знаки",
				newName: "??na??me.#1432&^$(@",
				wantErr: false,
			},
			{
				name:    "содержит только знаки",
				newName: "?>><#(*@$&",
				wantErr: false,
			},
			{
				name:    "содержит только пробелы",
				newName: " ",
				wantErr: true,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				input := UpdateNameInput{
					SubjectUserID: uuid.NewString(),
					ChatID:        uuid.NewString(),
					NewName:       tt.newName,
				}
				if err := input.Validate(); tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := UpdateNameInput{
			SubjectUserID: id,
			ChatID:        id,
			NewName:       "NewName",
		}
		return input.Validate()
	})
}

// Test_Chats_UpdateName тестирует обновления названия чата
func Test_Chats_UpdateName(t *testing.T) {
	t.Run("только существующий чат можно обновить", func(t *testing.T) {
		env := initChatTestEnv(t)
		input := UpdateNameInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
			NewName:       "newName",
		}
		// Обновить название чата
		chat, err := env.chatsService.UpdateName(input)
		// Вернется ошибка, потому что чата не существует
		assert.ErrorIs(t, err, ErrChatNotExists)
		assert.Zero(t, chat)
	})
	t.Run("только главный администратор может изменять название", func(t *testing.T) {
		env := initChatTestEnv(t)
		// Создать чат
		inputChatCreate := env.newCreateInputRandom()
		out, err := env.chatsService.Create(inputChatCreate)
		assert.NoError(t, err)
		assert.NotZero(t, out)
		// Попытаться изменить название от имени случайного пользователя
		input := UpdateNameInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        out.Chat.ID,
			NewName:       "newName",
		}
		updatedChat, err := env.chatsService.UpdateName(input)
		// Вернется ошибка, потому что пользователь не главный администратор чата
		assert.ErrorIs(t, err, ErrSubjectUserIsNotChief)
		assert.Zero(t, updatedChat)
	})
	t.Run("новое название чата сохранится и его можно прочитать", func(t *testing.T) {
		env := initChatTestEnv(t)
		// Создать чат
		inputChatCreate := env.newCreateInputRandom()
		out, err := env.chatsService.Create(inputChatCreate)
		assert.NoError(t, err)
		assert.NotZero(t, out)
		// Попытаться изменить название от имени случайного пользователя
		input := UpdateNameInput{
			SubjectUserID: out.Chat.ChiefUserID,
			ChatID:        out.Chat.ID,
			NewName:       "newName",
		}
		updatedChat, err := env.chatsService.UpdateName(input)
		assert.NoError(t, err)
		assert.NotZero(t, updatedChat)
		// Результат совпадает с входящими значениями
		assert.Equal(t, input.ChatID, updatedChat.ID)
		assert.Equal(t, input.NewName, updatedChat.Name)
		// Получить список чатов
		chats, err := env.chatsService.ChatsRepo.List(domain.ChatsFilter{})
		assert.NoError(t, err)
		// В списке этот чат будет единственным
		if assert.Len(t, chats, 1) {
			assert.Equal(t, input.ChatID, chats[0].ID)
			assert.Equal(t, input.NewName, chats[0].Name)
		}
	})
}
