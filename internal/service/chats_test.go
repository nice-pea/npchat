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
)

func (suite *servicesTestSuite) newUserChatsInput(userID string) WhichParticipateInput {
	return WhichParticipateInput{
		SubjectUserID: userID,
		UserID:        userID,
	}
}

func (suite *servicesTestSuite) assertChatEqualInput(in CreateInput, chat domain.Chat) {
	suite.Equal(in.Name, chat.Name)
	suite.Equal(in.ChiefUserID, chat.ChiefUserID)
}

func (suite *servicesTestSuite) newCreateInputRandom() CreateInput {
	return CreateInput{
		ChiefUserID: uuid.NewString(),
		Name:        fmt.Sprintf("name%d", rand.Int()),
	}
}

// Test_UserChatsInput_Validate тестирует валидацию входящих параметров запроса списка чатов в которых участвует пользователь
func Test_UserChatsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := WhichParticipateInput{
			SubjectUserID: id,
			UserID:        id,
		}
		return in.Validate()
	})
}

// Test_Chats_UserChats тестирует запрос список чатов в которых участвует пользователь
func (suite *servicesTestSuite) Test_Chats_UserChats() {
	suite.Run("пользователь может запрашивать только свой чат", func() {
		input := WhichParticipateInput{
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		chats, err := suite.ss.chats.UserChats(input)
		suite.ErrorIs(err, ErrUnauthorizedChatsView)
		suite.Empty(chats)
	})

	suite.Run("пустой список из пустого репозитория", func() {
		input := suite.newUserChatsInput(uuid.NewString())
		userChats, err := suite.ss.chats.UserChats(input)
		suite.NoError(err)
		suite.Empty(userChats)
	})

	suite.Run("пустой список если у пользователя нет чатов", func() {
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.saveChat(domain.Chat{
				ID: uuid.NewString(),
			})
			suite.saveMember(domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			})
		}
		input := suite.newUserChatsInput(uuid.NewString())
		userChats, err := suite.ss.chats.UserChats(input)
		suite.NoError(err)
		suite.Empty(userChats)
	})

	suite.Run("у пользователя может быть несколько чатов", func() {
		userID := uuid.NewString()
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.saveChat(domain.Chat{
				ID: uuid.NewString(),
			})
			suite.saveMember(domain.Member{
				ID:     uuid.NewString(),
				UserID: userID,
				ChatID: chat.ID,
			})
		}
		input := suite.newUserChatsInput(userID)
		userChats, err := suite.ss.chats.UserChats(input)
		suite.NoError(err)
		suite.Len(userChats, chatsAllCount)
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
func (suite *servicesTestSuite) Test_Chats_CreateChat() {
	suite.Run("выходящие совпадают с заданными", func() {
		// Создать чат
		input := suite.newCreateInputRandom()
		out, err := suite.ss.chats.Create(input)
		suite.NoError(err)
		// Сравнить результат с входящими значениями
		suite.Equal(input.ChiefUserID, out.ChiefMember.UserID)
		suite.assertChatEqualInput(input, out.Chat)
	})

	suite.Run("можно затем прочитать из репозитория", func() {
		// Создать чат
		input := suite.newCreateInputRandom()
		out, err := suite.ss.chats.Create(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// Получить список чатов
		chats, err := suite.ss.chats.ChatsRepo.List(domain.ChatsFilter{})
		suite.Require().NoError(err)
		// В списке этот чат будет единственным
		if suite.Len(chats, 1) {
			suite.assertChatEqualInput(input, chats[0])
		}
	})

	suite.Run("создается участник для главного администратора", func() {
		// Создать чат
		input := suite.newCreateInputRandom()
		out, err := suite.ss.chats.Create(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// Получить список участников
		members, err := suite.ss.chats.MembersRepo.List(domain.MembersFilter{})
		suite.NoError(err)
		// В списке этот участник будет единственным
		if suite.Len(members, 1) {
			// Участником является главный администратор созданного чата
			suite.Equal(input.ChiefUserID, members[0].UserID)
			suite.Equal(out.Chat.ID, members[0].ChatID)
		}
	})

	suite.Run("можно создать чаты с одинаковым именем", func() {
		input := suite.newCreateInputRandom()
		// Создать несколько чатов с одинаковым именем
		const chatsAllCount = 2
		for range chatsAllCount {
			out, err := suite.ss.chats.Create(input)
			suite.Require().NoError(err)
			suite.Require().NotZero(out)
		}
		// Получить список чатов
		chats, err := suite.ss.chats.ChatsRepo.List(domain.ChatsFilter{})
		suite.NoError(err)
		// Количество чатов равно количеству созданных
		suite.Len(chats, chatsAllCount)
	})

	suite.Run("количество созданных чатов на одного пользователя не ограничено", func() {
		// Пользователь
		userID := uuid.NewString()
		// Создать много чатов от лица пользователя
		const chatsAllCount = 900
		for range chatsAllCount {
			out, err := suite.ss.chats.Create(CreateInput{
				ChiefUserID: userID,
				Name:        "name",
			})
			suite.Require().NoError(err)
			suite.Require().NotZero(out)
		}
		// Получить список чатов
		chats, err := suite.ss.chats.ChatsRepo.List(domain.ChatsFilter{})
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
				input := UpdateNameInput{
					SubjectID: uuid.NewString(),
					ChatID:    uuid.NewString(),
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
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := UpdateNameInput{
			SubjectID: id,
			ChatID:    id,
			NewName:   "NewName",
		}
		return input.Validate()
	})
}

// Test_Chats_UpdateName тестирует обновления названия чата
func (suite *servicesTestSuite) Test_Chats_UpdateName() {
	suite.Run("только существующий чат можно обновить", func() {
		input := UpdateNameInput{
			SubjectID: uuid.NewString(),
			ChatID:    uuid.NewString(),
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
		createdOut, err := suite.ss.chats.Create(inputChatCreate)
		suite.Require().NoError(err)
		suite.Require().NotZero(createdOut)
		// Попытаться изменить название от имени случайного пользователя
		input := UpdateNameInput{
			SubjectID: uuid.NewString(),
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
		createdOut, err := suite.ss.chats.Create(inputChatCreate)
		suite.Require().NoError(err)
		suite.Require().NotZero(createdOut)
		// Попытаться изменить название от имени случайного пользователя
		input := UpdateNameInput{
			SubjectID: createdOut.Chat.ChiefUserID,
			ChatID:    createdOut.Chat.ID,
			NewName:   "newName",
		}
		updatedChat, err := suite.ss.chats.UpdateName(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(updatedChat)
		// Результат совпадает с входящими значениями
		suite.Require().Equal(input.ChatID, updatedChat.ID)
		suite.Require().Equal(input.NewName, updatedChat.Name)
		// Получить список чатов
		chats, err := suite.ss.chats.ChatsRepo.List(domain.ChatsFilter{})
		suite.Require().NoError(err)
		// В списке этот чат будет единственным
		if suite.Len(chats, 1) {
			suite.Equal(input.ChatID, chats[0].ID)
			suite.Equal(input.NewName, chats[0].Name)
		}
	})
}
