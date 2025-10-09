package deleteMember

import (
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

// Test_Members_DeleteMember тестирует удаление участника чата
func (suite *testSuite) Test_Members_DeleteMember() {
	suite.Run("нельзя удалить самого себя", func() {
		// Создать usecase и моки
		usecase, _, _ := newUsecase(suite)
		// Удалить участника
		userID := uuid.New()
		input := In{
			SubjectID: userID,
			ChatID:    uuid.New(),
			UserID:    userID,
		}
		out, err := usecase.DeleteMember(input)
		// Вернется ошибка, потому что пользователь пытается удалить самого себя
		suite.ErrorIs(err, ErrMemberCannotDeleteHimself)
		suite.Zero(out)
	})

	suite.Run("чат должен существовать", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		// Удалить участника
		input := In{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
			UserID:    uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return(nil, nil)
		out, err := usecase.DeleteMember(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Zero(out)
	})

	//suite.Run("subject должен быть участником чата", func() {
	//	// Создать чат
	//	chat := suite.upsertChat(suite.rndChat())
	//	// Удалить участника
	//	input := DeleteMemberIn{
	//		SubjectID: uuid.New(),
	//		ChatID:    chat.ID,
	//		UserID:    uuid.New(),
	//	}
	//	err := usecase.DeleteMember(input)
	//	// Вернется ошибка, потому что пользователь не участник чата
	//	suite.ErrorIs(err, ErrSubjectIsNotMember)
	//})

	suite.Run("subject должен быть главным администратором чата", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Удалить участника
		input := In{
			SubjectID: participant.UserID,
			ChatID:    chat.ID,
			UserID:    uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		out, err := usecase.DeleteMember(input)
		// Вернется ошибка, потому что участник не главный администратор
		suite.ErrorIs(err, ErrSubjectUserIsNotChief)
		suite.Zero(out)
	})

	suite.Run("user должен быть участником чата", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Удалить участника
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			UserID:    uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		out, err := usecase.DeleteMember(input)
		// Вернется ошибка, потому что удаляемый пользователь не является участником
		suite.ErrorIs(err, chatt.ErrParticipantNotExists)
		suite.Zero(out)
	})

	suite.Run("после удаления участник перестает быть участником", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return()
		// Создать чат
		chat := suite.RndChat()
		// Создать участника для удаления
		participant := suite.AddRndParticipant(&chat)
		// Удалить участника
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			UserID:    participant.UserID,
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil)
		out, err := usecase.DeleteMember(input)
		suite.Require().NoError(err)
		suite.Zero(out)
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		// Настройка мока
		var consumedEvents []events.Event
		mockEventsConsumer.EXPECT().Consume(mock.Anything).
			Run(func(events []events.Event) {
				consumedEvents = append(consumedEvents, events...)
			}).
			Return()

		// Создать чат
		chat := suite.RndChat()
		// Создать участника для удаления
		participant := suite.AddRndParticipant(&chat)
		// Удалить участника
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			UserID:    participant.UserID,
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil)
		out, err := usecase.DeleteMember(input)
		suite.Require().NoError(err)
		suite.Zero(out)

		// Проверить список опубликованных событий
		suite.AssertHasEventType(consumedEvents, chatt.EventParticipantRemoved)
	})
}

func newUsecase(suite *testSuite) (*DeleteMemberUsecase, *mockChatt.Repository, *mockEvents.Consumer) {
	uc := &DeleteMemberUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	mockEventsConsumer := uc.EventConsumer.(*mockEvents.Consumer)
	return uc, mockRepo, mockEventsConsumer
}
