package service

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/nice-pea/npchat/internal/domain/chatt"
)

func (suite *servicesTestSuite) newUserChatsInput(userID uuid.UUID) WhichParticipateIn {
	return WhichParticipateIn{
		SubjectID: userID,
		UserID:    userID,
	}
}

func (suite *servicesTestSuite) newCreateInputRandom() CreateChatIn {
	return CreateChatIn{
		ChiefUserID: uuid.New(),
		Name:        fmt.Sprintf("name%d", rand.Int()),
	}
}

// Test_Chats_UserChats тестирует запрос список чатов в которых участвует пользователь
func (suite *servicesTestSuite) Test_Chats_UserChats() {
	suite.Run("пользователь может запрашивать только свой чат", func() {
		input := WhichParticipateIn{
			SubjectID: uuid.New(),
			UserID:    uuid.New(),
		}
		chats, err := suite.ss.chats.WhichParticipate(input)
		suite.ErrorIs(err, ErrUnauthorizedChatsView)
		suite.Empty(chats)
	})

	suite.Run("пустой список из пустого репозитория", func() {
		input := suite.newUserChatsInput(uuid.New())
		userChats, err := suite.ss.chats.WhichParticipate(input)
		suite.NoError(err)
		suite.Empty(userChats)
	})

	suite.Run("пустой список если у пользователя нет чатов", func() {
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.upsertChat(suite.rndChat())
			suite.addRndParticipant(&chat)
		}
		input := suite.newUserChatsInput(uuid.New())
		userChats, err := suite.ss.chats.WhichParticipate(input)
		suite.NoError(err)
		suite.Empty(userChats)
	})

	suite.Run("у пользователя может быть несколько чатов", func() {
		userID := uuid.New()
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.upsertChat(suite.rndChat())
			p := suite.newParticipant(uuid.New())
			suite.addParticipant(&chat, p)
		}
		input := suite.newUserChatsInput(userID)
		userChats, err := suite.ss.chats.WhichParticipate(input)
		suite.NoError(err)
		suite.Len(userChats, chatsAllCount)
	})
}

// Test_Chats_CreateChat тестирует создание чата
func (suite *servicesTestSuite) Test_Chats_CreateChat() {
	suite.Run("выходящие совпадают с заданными", func() {
		// Создать чат
		input := suite.newCreateInputRandom()
		out, err := suite.ss.chats.CreateChat(input)
		suite.NoError(err)
		// Сравнить результат с входящими значениями
		suite.Equal(input.Name, out.Chat.Name)
		suite.Equal(input.ChiefUserID, out.Chat.ChiefID)
	})

	suite.Run("можно затем прочитать из репозитория", func() {
		// Создать чат
		input := suite.newCreateInputRandom()
		out, err := suite.ss.chats.CreateChat(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// Получить список чатов
		chats, err := suite.rr.chats.List(chatt.Filter{})
		suite.Require().NoError(err)
		// В списке этот чат будет единственным
		suite.Require().Len(chats, 1)
		suite.Equal(out.Chat.Name, chats[0].Name)
		suite.Equal(out.Chat.ChiefID, chats[0].ChiefID)
	})

	suite.Run("создается участник для главного администратора", func() {
		// Создать чат
		input := suite.newCreateInputRandom()
		out, err := suite.ss.chats.CreateChat(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// Получить список участников
		//members, err := suite.ss.chats.MembersRepo.List(domain.MembersFilter{})
		//suite.NoError(err)
		// В списке этот участник будет единственным
		suite.Require().Len(out.Chat.Participants, 1)
		// Участником является главный администратор созданного чата
		suite.Equal(input.ChiefUserID, out.Chat.Participants[0].UserID)
	})

	suite.Run("можно создать чаты с одинаковым именем", func() {
		input := suite.newCreateInputRandom()
		// Создать несколько чатов с одинаковым именем
		const chatsAllCount = 2
		for range chatsAllCount {
			out, err := suite.ss.chats.CreateChat(input)
			suite.Require().NoError(err)
			suite.Require().NotZero(out)
		}
		// Получить список чатов
		chats, err := suite.rr.chats.List(chatt.Filter{})
		suite.NoError(err)
		// Количество чатов равно количеству созданных
		suite.Len(chats, chatsAllCount)
	})

	suite.Run("количество созданных чатов на одного пользователя не ограничено", func() {
		// Пользователь
		userID := uuid.New()
		// Создать много чатов от лица пользователя
		const chatsAllCount = 900
		for range chatsAllCount {
			out, err := suite.ss.chats.CreateChat(CreateChatIn{
				ChiefUserID: userID,
				Name:        "name",
			})
			suite.Require().NoError(err)
			suite.Require().NotZero(out)
		}
		// Получить список чатов
		chats, err := suite.rr.chats.List(chatt.Filter{})
		suite.NoError(err)
		// Количество чатов равно количеству созданных
		suite.Len(chats, chatsAllCount)
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
				input := UpdateNameIn{
					SubjectID: uuid.New(),
					ChatID:    uuid.New(),
					NewName:   tt.newName,
				}
				if err := input.Validate(); tt.wantErr {
					assert.Error(t, err)
				} else {
					assert.NoError(t, err)
				}
			})
		}
	})
}

// Test_Chats_UpdateName тестирует обновления названия чата
func (suite *servicesTestSuite) Test_Chats_UpdateName() {
	suite.Run("только существующий чат можно обновить", func() {
		input := UpdateNameIn{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
			NewName:   "newName",
		}
		// Обновить название чата
		chat, err := suite.ss.chats.UpdateName(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
		suite.Zero(chat)
	})

	suite.Run("только главный администратор может изменять название", func() {
		// Создать чат
		inputChatCreate := suite.newCreateInputRandom()
		createdOut, err := suite.ss.chats.CreateChat(inputChatCreate)
		suite.Require().NoError(err)
		suite.Require().NotZero(createdOut)
		// Попытаться изменить название от имени случайного пользователя
		input := UpdateNameIn{
			SubjectID: uuid.New(),
			ChatID:    createdOut.Chat.ID,
			NewName:   "newName",
		}
		updatedChat, err := suite.ss.chats.UpdateName(input)
		// Вернется ошибка, потому что пользователь не главный администратор чата
		suite.ErrorIs(err, ErrSubjectUserIsNotChief)
		suite.Zero(updatedChat)
	})

	suite.Run("новое название чата сохранится и его можно прочитать", func() {
		// Создать чат
		inputChatCreate := suite.newCreateInputRandom()
		createdOut, err := suite.ss.chats.CreateChat(inputChatCreate)
		suite.Require().NoError(err)
		suite.Require().NotZero(createdOut)
		// Попытаться изменить название от имени случайного пользователя
		input := UpdateNameIn{
			SubjectID: createdOut.Chat.ChiefID,
			ChatID:    createdOut.Chat.ID,
			NewName:   "newName",
		}
		out, err := suite.ss.chats.UpdateName(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// Результат совпадает с входящими значениями
		suite.Require().Equal(input.ChatID, out.Chat.ID)
		suite.Require().Equal(input.NewName, out.Chat.Name)
		// Получить список чатов
		chats, err := suite.rr.chats.List(chatt.Filter{})
		suite.Require().NoError(err)
		// В списке этот чат будет единственным
		if suite.Len(chats, 1) {
			suite.Equal(input.ChatID, chats[0].ID)
			suite.Equal(input.NewName, chats[0].Name)
		}
	})
}
