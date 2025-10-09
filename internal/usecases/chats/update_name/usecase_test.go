package updateName

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	mockChatt "github.com/nice-pea/npchat/internal/domain/chatt/mocks"
	"github.com/nice-pea/npchat/internal/usecases/events"
	mockEvents "github.com/nice-pea/npchat/internal/usecases/events/mocks"
	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
)

type testSuite struct {
	serviceSuite.SuiteWithMocks
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

// Test_Chats_UpdateName тестирует обновления названия чата
func (suite *testSuite) Test_Chats_UpdateName() {

	suite.Run("только существующий чат можно обновить", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		input := In{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
			NewName:   "newName",
		}
		// Обновить название чата
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{}, nil)
		chat, err := usecase.UpdateName(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Zero(chat)
	})

	suite.Run("только главный администратор может изменять название", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Попытаться изменить название от имени случайного пользователя
		input := In{
			SubjectID: uuid.New(),
			ChatID:    chat.ID,
			NewName:   "newName",
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		updatedChat, err := usecase.UpdateName(input)
		// Вернется ошибка, потому что пользователь не главный администратор чата
		suite.ErrorIs(err, ErrSubjectUserIsNotChief)
		suite.Zero(updatedChat)
	})

	suite.Run("новое название чата сохранится и его можно прочитать", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventConsumer := newUsecase(suite)
		mockEventConsumer.EXPECT().Consume(mock.Anything).Return().Once()
		// Создать чат
		chat := suite.RndChat()
		// Изменить название от имени администратора
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			NewName:   "newName",
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		mockRepo.EXPECT().Upsert(mock.Anything).Run(func(chat chatt.Chat) {
			suite.Equal(input.NewName, chat.Name)
		}).Return(nil)
		out, err := usecase.UpdateName(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// Результат совпадает с входящими значениями
		suite.Require().Equal(input.ChatID, out.Chat.ID)
		suite.Require().Equal(input.NewName, out.Chat.Name)
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventConsumer := newUsecase(suite)
		// Настройка мока
		var consumedEvents []events.Event
		mockEventConsumer.EXPECT().Consume(mock.Anything).Run(func(events []events.Event) {
			consumedEvents = append(consumedEvents, events...)
		}).Return()
		// Создать чат
		chat := suite.RndChat()
		// Изменить название от имени администратора
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			NewName:   "newName",
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil)
		out, err := usecase.UpdateName(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Проверить список опубликованных событий
		suite.AssertHasEventType(consumedEvents, chatt.EventChatUpdated)
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

func newUsecase(suite *testSuite) (*UpdateNameUsecase, *mockChatt.Repository, *mockEvents.Consumer) {
	uc := &UpdateNameUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	mockEventsConsumer := uc.EventConsumer.(*mockEvents.Consumer)
	return uc, mockRepo, mockEventsConsumer
}
