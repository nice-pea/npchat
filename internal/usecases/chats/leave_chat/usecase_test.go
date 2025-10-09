package leaveChat

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
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

// Test_Members_LeaveChat тестирует выход участника из чата
func (suite *testSuite) Test_Members_LeaveChat() {
	suite.Run("чат должен существовать", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		// Покинуть чат
		input := In{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return(nil, nil).Once()
		out, err := usecase.LeaveChat(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Zero(out)
	})

	suite.Run("пользователь должен быть участником чата", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Покинуть чат
		input := In{
			SubjectID: uuid.New(),
			ChatID:    chat.ID,
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()
		out, err := usecase.LeaveChat(input)
		// Вернется ошибка, потому что пользователь не участник чата
		suite.ErrorIs(err, chatt.ErrParticipantNotExists)
		suite.Zero(out)
	})

	suite.Run("пользователь не должен быть главным администратором чата", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		chat := suite.RndChat()
		// Покинуть чат
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()
		out, err := usecase.LeaveChat(input)
		// Вернется ошибка, потому что пользователь главный администратор чата
		suite.ErrorIs(err, chatt.ErrCannotRemoveChief)
		suite.Zero(out)
	})

	suite.Run("после выхода пользователь перестает быть участником", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Once()
		// Создать чат
		chat := suite.RndChat()
		// Создать участника в этом чате
		participant := suite.AddRndParticipant(&chat)
		// Покинуть чат
		input := In{
			SubjectID: participant.UserID,
			ChatID:    chat.ID,
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()
		mockRepo.EXPECT().Upsert(mock.Anything).Run(func(chatRepo chatt.Chat) {
			chat = chatRepo
		}).Return(nil)
		out, err := usecase.LeaveChat(input)
		suite.Require().NoError(err)
		suite.Zero(out)
		suite.Len(chat.Participants, 1)
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Новый экземпляр usecase
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		// Настройка мока
		var consumedEvents []events.Event
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Run(func(events []events.Event) {
			consumedEvents = append(consumedEvents, events...)
		}).Return().Once()
		// Создать чат
		chat := suite.RndChat()
		// Создать участника в этом чате
		participant := suite.AddRndParticipant(&chat)
		// Покинуть чат
		input := In{
			SubjectID: participant.UserID,
			ChatID:    chat.ID,
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()
		chat.Participants = chat.Participants[:1]
		mockRepo.EXPECT().Upsert(chat).Return(nil).Once()
		out, err := usecase.LeaveChat(input)
		suite.Require().NoError(err)
		suite.Zero(out)

		// Проверить список опубликованных событий
		suite.AssertHasEventType(consumedEvents, chatt.EventParticipantRemoved)
	})

	suite.Run("отправленные приглашения участника, отменятся вместе с его выходом", func() {
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)

		var consumedEvents []events.Event
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Run(func(e []events.Event) {
			consumedEvents = append(consumedEvents, e...)
		}).Return().Once()

		chat := suite.RndChat()
		participant := suite.AddRndParticipant(&chat)

		invitation := suite.NewInvitation(participant.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)

		input := In{
			SubjectID: participant.UserID,
			ChatID:    chat.ID,
		}

		chatBeforeLeave := chat
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chatBeforeLeave}, nil).Once()
		mockRepo.EXPECT().Upsert(mock.Anything).Run(func(updated chatt.Chat) {
			suite.False(updated.HasParticipant(participant.UserID))
			suite.Empty(updated.Invitations)
		}).Return(nil).Once()

		out, err := usecase.LeaveChat(input)
		suite.Require().NoError(err)
		suite.Zero(out)

		suite.AssertHasEventType(consumedEvents, chatt.EventParticipantRemoved)
		suite.AssertHasEventType(consumedEvents, chatt.EventInvitationRemoved)
	})
}

func newUsecase(suite *testSuite) (*LeaveChatUsecase, *mockChatt.Repository, *mockEvents.Consumer) {
	uc := &LeaveChatUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	mockEventsConsumer := uc.EventConsumer.(*mockEvents.Consumer)
	return uc, mockRepo, mockEventsConsumer
}
