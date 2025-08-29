package updateName

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	createChat "github.com/nice-pea/npchat/internal/service/chats/create_chat"
	serviceSuite "github.com/nice-pea/npchat/internal/service/suite"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

func (suite *testSuite) newCreateInputRandom() createChat.In {
	return createChat.In{
		ChiefUserID: uuid.New(),
		Name:        fmt.Sprintf("name%d", rand.Int()),
	}
}

// Test_Chats_UpdateName тестирует обновления названия чата
func (suite *testSuite) Test_Chats_UpdateName() {
	usecase := &UpdateNameUsecase{
		Repo: suite.RR.Chats,
	}
	createChatUsecase := createChat.CreateChatUsecase{
		Repo: suite.RR.Chats,
	}

	suite.Run("только существующий чат можно обновить", func() {
		input := In{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
			NewName:   "newName",
		}
		// Обновить название чата
		chat, err := usecase.UpdateName(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Zero(chat)
	})

	suite.Run("только главный администратор может изменять название", func() {
		// Создать чат
		inputChatCreate := suite.newCreateInputRandom()
		createdOut, err := createChatUsecase.CreateChat(inputChatCreate)
		suite.Require().NoError(err)
		suite.Require().NotZero(createdOut)
		// Попытаться изменить название от имени случайного пользователя
		input := In{
			SubjectID: uuid.New(),
			ChatID:    createdOut.Chat.ID,
			NewName:   "newName",
		}
		updatedChat, err := usecase.UpdateName(input)
		// Вернется ошибка, потому что пользователь не главный администратор чата
		suite.ErrorIs(err, ErrSubjectUserIsNotChief)
		suite.Zero(updatedChat)
	})

	suite.Run("новое название чата сохранится и его можно прочитать", func() {
		// Создать чат
		inputChatCreate := suite.newCreateInputRandom()
		createdOut, err := createChatUsecase.CreateChat(inputChatCreate)
		suite.Require().NoError(err)
		suite.Require().NotZero(createdOut)
		// Попытаться изменить название от имени случайного пользователя
		input := In{
			SubjectID: createdOut.Chat.ChiefID,
			ChatID:    createdOut.Chat.ID,
			NewName:   "newName",
		}
		out, err := usecase.UpdateName(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// Результат совпадает с входящими значениями
		suite.Require().Equal(input.ChatID, out.Chat.ID)
		suite.Require().Equal(input.NewName, out.Chat.Name)
		// Получить список чатов
		chats, err := suite.RR.Chats.List(chatt.Filter{})
		suite.Require().NoError(err)
		// В списке этот чат будет единственным
		if suite.Len(chats, 1) {
			suite.Equal(input.ChatID, chats[0].ID)
			suite.Equal(input.NewName, chats[0].Name)
		}
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
				input := In{
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
