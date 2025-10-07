package acceptInvitation

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

// Test_Invitations_AcceptInvitation тестирует принятие приглашения
func (suite *testSuite) Test_Invitations_AcceptInvitation() {

	suite.Run("приглашение должно существовать", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEvent := newUsecase(suite)
		// Настройка мока
		mockEvent.EXPECT().Consume(mock.Anything).Return().Maybe()

		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		p := suite.AddRndParticipant(&chat)
		// Сохранить чат
		// настроить мок
		mockRepo.EXPECT().Upsert(chat).Return(nil)
		suite.UpsertChat(chat)
		// Принять приглашение
		input := In{
			SubjectID:    p.UserID,
			InvitationID: uuid.New(),
		}
		// настроить мок
		mockRepo.EXPECT().List(chatt.Filter{
			InvitationID: input.InvitationID,
		}).Return(nil, chatt.ErrChatNotExists)
		// suite.SetupAcceptInvitationMocks(input.InvitationID, chat)
		out, err := usecase.AcceptInvitation(input)
		suite.ErrorIs(err, ErrInvitationNotExists)
		suite.Zero(out)
	})

	suite.Run("приняв приглашение, пользователь становится участником чата", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEvent := newUsecase(suite)
		// Настройка мока
		mockEvent.EXPECT().Consume(mock.Anything).Return().Maybe()
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		p := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(p.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Сохранить чат
		// настроить мок
		mockRepo.EXPECT().Upsert(chat).Return(nil)
		suite.UpsertChat(chat)
		// Принять приглашение
		input := In{
			SubjectID:    invitation.RecipientID,
			InvitationID: invitation.ID,
		}
		// настройка моков
		suite.SetupAcceptInvitationMocks(input.InvitationID, chat)
		out, err := usecase.AcceptInvitation(input)
		suite.Zero(out)
		suite.Require().NoError(err)
		// Получить список участников
		mockRepo.EXPECT().List(chatt.Filter{}).RunAndReturn(func(filter chatt.Filter) ([]chatt.Chat, error) {
			chat.Participants = append(chat.Participants, chatt.Participant{
				UserID: input.InvitationID,
			})
			return []chatt.Chat{chat}, nil
		}).Once()
		chats, err := suite.RR.Chats.List(chatt.Filter{})
		suite.NoError(err)
		// В списке будет 3 участника: адм., приглашаемый, приглашающий
		suite.Require().Len(chats, 1)
		suite.Require().Len(chats[0].Participants, 3)
		suite.Contains(chats[0].Participants, p)
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEvent := newUsecase(suite)
		// Настройка мока
		var consumedEvents []events.Event
		mockEvent.EXPECT().Consume(mock.Anything).Run(
			func(args []events.Event) {
				consumedEvents = append(consumedEvents, args...)
			}).Return()
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		p := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(p.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Сохранить чат
		mockRepo.EXPECT().Upsert(chat).Return(nil)
		suite.UpsertChat(chat)
		// Принять приглашение
		input := In{
			SubjectID:    invitation.RecipientID,
			InvitationID: invitation.ID,
		}
		// настройка моков
		suite.SetupAcceptInvitationMocks(input.InvitationID, chat)
		out, err := usecase.AcceptInvitation(input)
		suite.Zero(out)
		suite.Require().NoError(err)

		// Проверить список опубликованных событий
		suite.AssertHasEventType(consumedEvents, chatt.EventInvitationRemoved)
		suite.AssertHasEventType(consumedEvents, chatt.EventParticipantAdded)
	})
}

func newUsecase(suite *testSuite) (*AcceptInvitationUsecase, *mockChatt.Repository, *mockEvents.Consumer) {
	uc := &AcceptInvitationUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	mockEvent := uc.EventConsumer.(*mockEvents.Consumer)
	return uc, mockRepo, mockEvent
}
