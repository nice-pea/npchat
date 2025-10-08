package sendInvitation

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

// Test_Invitations_SendChatInvitation тестирует отправку приглашения
func (suite *testSuite) Test_Invitations_SendChatInvitation() {

	suite.Run("чат должен существовать", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Maybe()
		// Отправить приглашение
		input := In{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
			UserID:    uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return(nil, nil)
		invitation, err := usecase.SendInvitation(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Zero(invitation)
	})

	suite.Run("субъект должен быть участником", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Maybe()
		// Создать чат
		chat := suite.RndChat()
		// Отправить приглашение
		input := In{
			SubjectID: uuid.New(),
			ChatID:    chat.ID,
			UserID:    uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		invitation, err := usecase.SendInvitation(input)
		// Вернется ошибка, потому что субъект не является участником чата
		suite.ErrorIs(err, chatt.ErrSubjectIsNotMember)
		suite.Zero(invitation)
	})

	suite.Run("приглашаемый пользователь может не существовать", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Maybe()
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Отправить приглашение
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		invitaion := suite.NewInvitation(input.SubjectID, input.UserID)
		mockRepo.EXPECT().Upsert(mock.Anything).Run(func(chat chatt.Chat) {
			suite.Equal(invitaion.RecipientID, input.UserID)
			suite.Equal(invitaion.SubjectID, input.SubjectID)
		}).Return(nil)
		out, err := usecase.SendInvitation(input)
		suite.NoError(err)
		suite.NotZero(out)
	})

	suite.Run("приглашаемый пользователь не должен состоять в этом чате", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Maybe()
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать участника для приглашаемого пользователя
		participantInvitating := suite.AddRndParticipant(&chat)
		// Отправить приглашение
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    participantInvitating.UserID,
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil)
		invitation, err := usecase.SendInvitation(input)
		// Вернется ошибка, потому что приглашаемый пользователь уже является участником этого чата
		suite.ErrorIs(err, chatt.ErrParticipantExists)
		suite.Zero(invitation)
	})

	suite.Run("одновременно не может существовать несколько приглашений одного пользователя в этот чат", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Maybe()
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать приглашаемого пользователя
		targetUserID := uuid.New()
		// Отправить приглашение
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    targetUserID,
		}

		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		invitation, err := usecase.SendInvitation(input)
		suite.NoError(err)
		suite.Require().NotZero(invitation)
		// Отправить повторно приглашение
		chat.Invitations = append(chat.Invitations, invitation.Invitation)
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()
		mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
		invitation, err = usecase.SendInvitation(input)
		// Вернется ошибка, потому что этот пользователь уже приглашен в чат
		suite.ErrorIs(err, chatt.ErrUserIsAlreadyInvited)
		suite.Zero(invitation)
	})

	suite.Run("любой участник может приглашать много пользователей", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Return().Maybe()
		// Создать чат
		chat := suite.RndChat()
		// Создать много приглашений от разных участников
		for range 5 {
			// Создать участника
			participant := suite.AddRndParticipant(&chat)
			for range 5 {
				// Отправить приглашение
				input := In{
					ChatID:    chat.ID,
					SubjectID: participant.UserID,
					UserID:    uuid.New(),
				}
				mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()
				mockRepo.EXPECT().Upsert(mock.Anything).Return(nil).Once()
				out, err := usecase.SendInvitation(input)
				suite.NoError(err)
				suite.Require().NotZero(out)
			}
		}
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Создать usecase и моки
		usecase, mockRepo, mockEventsConsumer := newUsecase(suite)
		// Настройка мока
		var consumedEvents []events.Event
		mockEventsConsumer.EXPECT().Consume(mock.Anything).Run(func(events []events.Event) {
			consumedEvents = append(consumedEvents, events...)
		}).Return()

		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Отправить приглашение
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()
		mockRepo.EXPECT().Upsert(mock.Anything).Run(func(chat chatt.Chat) {
			suite.Equal(input.UserID, chat.Invitations[0].RecipientID)
			suite.Equal(input.SubjectID, chat.Invitations[0].SubjectID)
		}).Return(nil)
		out, err := usecase.SendInvitation(input)
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить список опубликованных событий
		suite.AssertHasEventType(consumedEvents, chatt.EventInvitationAdded)
	})
}

func newUsecase(suite *testSuite) (*SendInvitationUsecase, *mockChatt.Repository, *mockEvents.Consumer) {
	uc := &SendInvitationUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	mockEventsConsumer := uc.EventConsumer.(*mockEvents.Consumer)
	return uc, mockRepo, mockEventsConsumer
}
