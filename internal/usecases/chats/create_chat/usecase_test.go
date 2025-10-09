package createChat

import (
	"fmt"
	"math/rand/v2"
	"testing"

	"github.com/google/uuid"
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

func (suite *testSuite) newCreateInputRandom() In {
	return In{
		ChiefUserID: uuid.New(),
		Name:        fmt.Sprintf("name%d", rand.Int()),
	}
}

// Test_Chats_CreateChat тестирует создание чата
func (suite *testSuite) Test_Chats_CreateChat() {

	suite.Run("выходящие совпадают с заданными", func() {
		usecase, mockRepo, mockEventConsumer := newUsecase(suite)
		mockEventConsumer.EXPECT().Consume(mock.Anything).Return().Once()
		// Создать чат
		input := suite.newCreateInputRandom()
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		out, err := usecase.CreateChat(input)
		suite.NoError(err)
		// Сравнить результат с входящими значениями
		suite.Equal(input.Name, out.Chat.Name)
		suite.Equal(input.ChiefUserID, out.Chat.ChiefID)
	})

	suite.Run("можно затем прочитать из репозитория", func() {
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		// Настройка мока
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Once()
		// Создать чат
		input := suite.newCreateInputRandom()
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		out, err := usecase.CreateChat(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
	})

	suite.Run("создается участник для главного администратора", func() {
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		// Настройка мока
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Once()
		// Создать чат
		input := suite.newCreateInputRandom()
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		out, err := usecase.CreateChat(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// В списке этот участник будет единственным
		suite.Require().Len(out.Chat.Participants, 1)
		// Участником является главный администратор созданного чата
		suite.Equal(input.ChiefUserID, out.Chat.Participants[0].UserID)
	})

	suite.Run("можно создать чаты с одинаковым именем", func() {
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		// Настройка мока
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return()
		input := suite.newCreateInputRandom()
		// Создать несколько чатов с одинаковым именем
		const chatsAllCount = 2
		chats := make([]chatt.Chat, chatsAllCount)
		for i := range chatsAllCount {
			mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
			out, err := usecase.CreateChat(input)
			suite.Require().NoError(err)
			suite.Require().NotZero(out)
			chats[i] = out.Chat
		}
	})

	suite.Run("количество созданных чатов на одного пользователя не ограничено", func() {
		const chatsAllCount = 900
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		// Настройка мока
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Times(chatsAllCount)
		// Пользователь
		userID := uuid.New()
		// Создать много чатов от лица пользователя
		for range chatsAllCount {
			mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
			out, err := usecase.CreateChat(In{
				ChiefUserID: userID,
				Name:        "name",
			})
			suite.Require().NoError(err)
			suite.Require().NotZero(out)
		}
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Новый экземпляр usecase
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		// Настройка мока
		var consumedEvents []events.Event
		mockEventsConsumer.EXPECT().Consume(mock.Anything).
			Run(func(events []events.Event) {
				consumedEvents = append(consumedEvents, events...)
			}).
			Return().Once()

		// Создать чат
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		out, err := usecase.CreateChat(In{
			ChiefUserID: uuid.New(),
			Name:        "name",
		})
		suite.Require().NoError(err)
		suite.Require().NotZero(out)

		// Проверить список опубликованных событий
		suite.AssertHasEventType(consumedEvents, chatt.EventChatCreated)
	})
}

func newUsecase(suite *testSuite) (*CreateChatUsecase, *mockChatt.Repository, *mockEvents.Consumer) {
	uc := &CreateChatUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	mockEventsConsumer := uc.EventConsumer.(*mockEvents.Consumer)
	return uc, mockRepo, mockEventsConsumer
}
