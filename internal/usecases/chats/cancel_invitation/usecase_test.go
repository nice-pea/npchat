package cancelInvitation

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

// Test_Invitations_CancelInvitation тестирует отмену приглашения
func (suite *testSuite) Test_Invitations_CancelInvitation() {
	suite.Run("приглашение должно существовать", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		// Отменить приглашение
		input := In{
			SubjectID:    uuid.New(),
			InvitationID: uuid.New(),
		}
		// Настройка мока
		mockRepo.EXPECT().
			List(mock.Anything).
			Return([]chatt.Chat{}, chatt.ErrChatNotExists).Once()
		out, err := usecase.CancelInvitation(input)
		// Вернется ошибка, потому что приглашения не существует
		suite.ErrorIs(err, ErrInvitationNotExists)
		suite.Zero(out)
	})

	suite.Run("приглашение могут отменить только пригласивший и приглашаемый пользователи, и администратор чата", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Times(3)
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Объявить id приглашаемого пользователя
		recipientID := uuid.New()
		// Список id тех пользователей, которые могут отменять приглашение
		cancelInvitationSubjectIDs := []uuid.UUID{
			chat.ChiefID,       // главный администратор
			participant.UserID, // пригласивший
			recipientID,        // приглашаемый
		}
		// Каждый попытается отменить приглашение
		for _, subjectUserID := range cancelInvitationSubjectIDs {
			// Создать приглашение
			invitation := suite.NewInvitation(participant.UserID, recipientID)
			suite.AddInvitation(&chat, invitation)
			// Отменить приглашение
			input := In{
				SubjectID:    subjectUserID,
				InvitationID: invitation.ID,
			}
			mockRepo.EXPECT().
				List(mock.Anything).
				Return([]chatt.Chat{chat}, nil).Once()
			mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
			out, err := usecase.CancelInvitation(input)
			suite.NoError(err)
			suite.Zero(out)
		}
	})

	suite.Run("другие участники не могут отменять приглашать ", func() {
		// Создать usecase и моки
		usecase, mockRepo, _ := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Случайный участник
		participantOther := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(participant.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Отменить приглашение
		input := In{
			SubjectID:    participantOther.UserID,
			InvitationID: invitation.ID,
		}
		mockRepo.EXPECT().
			List(mock.Anything).
			Return([]chatt.Chat{chat}, nil).Once()
		mockRepo.EXPECT().Upsert(chat).Return(nil).Once()
		out, err := usecase.CancelInvitation(input)
		// Вернется ошибка, потому что случайный участник не может отменять приглашение
		suite.ErrorIs(err, ErrSubjectUserNotAllowed)
		suite.Zero(out)
	})

	suite.Run("после отмены, приглашение удаляется", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return()
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(participant.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Отменить приглашение
		input := In{
			SubjectID:    invitation.SubjectID,
			InvitationID: invitation.ID,
		}
		mockRepo.EXPECT().
			List(mock.Anything).
			Return([]chatt.Chat{chat}, nil).Once()
		chat.Invitations = []chatt.Invitation{}
		mockRepo.EXPECT().Upsert(chat).RunAndReturn(func(chat chatt.Chat) error {
			suite.Empty(chat.Invitations)
			return nil
		}).Once()
		out, err := usecase.CancelInvitation(input)
		suite.Require().NoError(err)
		suite.Zero(out)
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		// Настройка мока
		var consumedEvents []events.Event
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Run(func(events []events.Event) {
			consumedEvents = append(consumedEvents, events...)
		}).Return().Once()

		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(participant.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)

		// Отменить приглашение
		input := In{
			SubjectID:    invitation.SubjectID,
			InvitationID: invitation.ID,
		}
		mockRepo.EXPECT().
			List(mock.Anything).
			Return([]chatt.Chat{chat}, nil).Once()
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		out, err := usecase.CancelInvitation(input)
		suite.Require().NoError(err)
		suite.Zero(out)

		// Проверить список опубликованных событий
		suite.AssertHasEventType(consumedEvents, chatt.EventInvitationRemoved)
	})
}

func newUsecase(suite *testSuite) (*CancelInvitationUsecase, *mockChatt.Repository, *mockEvents.Consumer) {
	uc := &CancelInvitationUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	mockEventsConsumer := uc.EventConsumer.(*mockEvents.Consumer)
	return uc, mockRepo, mockEventsConsumer
}
