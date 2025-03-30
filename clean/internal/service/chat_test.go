package service

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite/memory"
)

// newChatsService создает объект сервиса Chats с sqlite/memory репозиториями
func newChatsService(t *testing.T) *Chats {
	sqLiteInMemory, err := memory.Init(memory.Config{MigrationsDir: "../../migrations/repository/sqlite/memory"})
	assert.NoError(t, err)
	chatsRepository, err := sqLiteInMemory.NewChatsRepository()
	assert.NoError(t, err)
	membersRepository, err := sqLiteInMemory.NewMembersRepository()
	assert.NoError(t, err)
	return &Chats{
		ChatsRepo:   chatsRepository,
		MembersRepo: membersRepository,
	}
}

// Test_UserChatsInput_Validate тестирует валидацию входящих параметров запроса списка чатов в которых участвует пользователь
func Test_UserChatsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := UserChatsInput{
			SubjectUserID: id,
			UserID:        uuid.NewString(),
		}
		return in.Validate()
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := UserChatsInput{
			SubjectUserID: uuid.NewString(),
			UserID:        id,
		}
		return in.Validate()
	})
}

// Test_Chats_UserChats тестирует запрос список чатов в которых участвует пользователь
func Test_Chats_UserChats(t *testing.T) {
	t.Run("пустой список из пустого репозитория", func(t *testing.T) {
		id := uuid.NewString()
		input := UserChatsInput{
			SubjectUserID: id,
			UserID:        id,
		}
		userChats, err := newChatsService(t).UserChats(input)
		assert.NoError(t, err)
		assert.Len(t, userChats, 0)
	})
	t.Run("пустой список если у пользователя нет чатов", func(t *testing.T) {
		chatsService := newChatsService(t)
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
		input := UserChatsInput{
			SubjectUserID: userID,
			UserID:        userID,
		}
		userChats, err := chatsService.UserChats(input)
		assert.NoError(t, err)
		assert.Len(t, userChats, 0)
	})
	t.Run("UserID и SubjectUserID должны быть одинаковыми", func(t *testing.T) {
		chatsService := newChatsService(t)
		input := UserChatsInput{
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		chats, err := chatsService.UserChats(input)
		assert.Error(t, err)
		assert.Len(t, chats, 0)
	})
	t.Run("у пользователя может быть несколько чатов", func(t *testing.T) {
		chatsService := newChatsService(t)
		var existsChats []domain.Chat
		const countAll = 10
		for i := range [countAll]int{} {
			// Создать чат
			existsChats = append(existsChats, domain.Chat{
				ID:          uuid.NewString(),
				Name:        "name" + string(rune(i)),
				ChiefUserID: uuid.NewString(),
			})
			err := chatsService.ChatsRepo.Save(existsChats[i])
			assert.NoError(t, err)
		}
		userID := uuid.NewString()
		const count = countAll / 3
		for i := range count {
			err := chatsService.MembersRepo.Save(domain.Member{
				ID:     uuid.NewString(),
				UserID: userID,
				ChatID: existsChats[i].ID,
			})
			assert.NoError(t, err)
		}
		input := UserChatsInput{
			SubjectUserID: userID,
			UserID:        userID,
		}
		userChats, err := chatsService.UserChats(input)
		assert.NoError(t, err)
		assert.Len(t, userChats, 3)
	})
}

// Test_CreateChatInput_Validate тестирует валидацию входящих параметров для создания чата
func Test_CreateChatInput_Validate(t *testing.T) {
	t.Run("ошибка при пустом name", func(t *testing.T) {
		input := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "",
		}
		assert.Error(t, input.Validate())
	})
	t.Run("ошибка при пустом ChiefUserID", func(t *testing.T) {
		input := CreateInput{
			ChiefUserID: "",
			Name:        "qf",
		}
		assert.Error(t, input.Validate())
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := CreateInput{
			ChiefUserID: id,
			Name:        "qf",
		}
		return in.Validate()
	})
}

// Test_Chats_CreateChat тестирует создание чата
func Test_Chats_CreateChat(t *testing.T) {
	assertChatEqualIn := func(in CreateInput, out domain.Chat) {
		assert.Equal(t, out.ChiefUserID, in.ChiefUserID)
		assert.Equal(t, out.Name, in.Name)
	}
	t.Run("создание чата без ошибок", func(t *testing.T) {
		input := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "Name",
		}
		out, err := newChatsService(t).Create(input)
		assert.NoError(t, err)
		assert.NotZero(t, out)
	})
	t.Run("выходящие совпадают с заданными", func(t *testing.T) {
		input := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "Name",
		}
		out, err := newChatsService(t).Create(input)
		assert.NoError(t, err)
		assert.NotZero(t, out)
		assertChatEqualIn(input, out.Chat)
	})
	t.Run("возвращается чат с новым id", func(t *testing.T) {
		input := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "Name",
		}
		out, err := newChatsService(t).Create(input)
		assert.NoError(t, err)
		assert.NotZero(t, out)
		assert.NotZero(t, out.Chat.ID)
	})
	t.Run("можно затем прочитать из репозитория", func(t *testing.T) {
		chatsService := newChatsService(t)
		input := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "Name",
		}
		out, err := chatsService.Create(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, out.Chat.ID)
		chats, err := chatsService.ChatsRepo.List(domain.ChatsFilter{})
		assert.NoError(t, err)
		if assert.Len(t, chats, 1) {
			assertChatEqualIn(input, chats[0])
		}
	})
	t.Run("создается участник для главного администратора", func(t *testing.T) {
		chatsService := newChatsService(t)
		input := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "Name",
		}
		out, err := chatsService.Create(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, out.Chat.ID)
		assertChatEqualIn(input, out.Chat)
		members, err := chatsService.MembersRepo.List(domain.MembersFilter{})
		assert.NoError(t, err)
		if assert.Len(t, members, 1) {
			assert.Equal(t, input.ChiefUserID, members[0].UserID)
			assert.Equal(t, out.Chat.ID, members[0].ChatID)
		}
	})
	t.Run("пользователь создал много чатов", func(t *testing.T) {
		chatsService := newChatsService(t)
		userID := uuid.NewString()
		count := 900
		for i := 0; i < count; i++ {
			input := CreateInput{
				ChiefUserID: userID,
				Name:        fmt.Sprintf("Name%d", i),
			}
			out, err := chatsService.Create(input)
			assert.NoError(t, err)
			assert.NotZero(t, out)
		}
		list, err := chatsService.ChatsRepo.List(domain.ChatsFilter{})
		assert.NoError(t, err)
		assert.Len(t, list, count)
	})
	t.Run("можно создать с одинаковыми параметрами", func(t *testing.T) {
		chatsService := newChatsService(t)
		const count = 20
		input := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "name",
		}
		for range count {
			out, err := chatsService.Create(input)
			assert.NoError(t, err)
			assert.NotZero(t, out)
			assertChatEqualIn(input, out.Chat)
		}
		chats, err := chatsService.ChatsRepo.List(domain.ChatsFilter{})
		assert.NoError(t, err)
		assert.Len(t, chats, count)
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
			ChatID:        uuid.NewString(),
			NewName:       "NewName",
		}
		return input.Validate()
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := UpdateNameInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        id,
			NewName:       "NewName",
		}
		return input.Validate()
	})
}

// Test_Chats_UpdateName тестирует обновления названия чата
func Test_Chats_UpdateName(t *testing.T) {
	t.Run("без ошибок", func(t *testing.T) {
		chatsService := newChatsService(t)
		createInput := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "oldName",
		}
		createOut, err := chatsService.Create(createInput)
		assert.NoError(t, err)
		assert.NotZero(t, createOut)
		input := UpdateNameInput{
			SubjectUserID: createInput.ChiefUserID,
			ChatID:        createOut.Chat.ID,
			NewName:       "newName",
		}
		chat, err := chatsService.UpdateName(input)
		assert.NoError(t, err)
		assert.NotZero(t, chat)
		chats, err := chatsService.ChatsRepo.List(domain.ChatsFilter{})
		assert.NoError(t, err)
		if assert.Len(t, chats, 1) {
			assert.Equal(t, input.NewName, chats[0].Name)
		}
	})
	t.Run("только существующий чат можно обновить", func(t *testing.T) {
		chatsService := newChatsService(t)
		input := UpdateNameInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
			NewName:       "newName",
		}
		chat, err := chatsService.UpdateName(input)
		assert.Error(t, err)
		assert.Zero(t, chat)
	})
	t.Run("только chief может изменять название", func(t *testing.T) {
		chatsService := newChatsService(t)
		createInput := CreateInput{
			ChiefUserID: uuid.NewString(),
			Name:        "oldName",
		}
		createOut, err := chatsService.Create(createInput)
		assert.NoError(t, err)
		assert.NotZero(t, createOut)
		input := UpdateNameInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        createOut.Chat.ID,
			NewName:       "newName",
		}
		updatedChat, err := chatsService.UpdateName(input)
		assert.Error(t, err)
		assert.Zero(t, updatedChat)
	})
}
