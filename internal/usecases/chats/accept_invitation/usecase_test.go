package acceptInvitation

import (
	"testing"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/mock"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/chatt"
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

// Test_Invitations_AcceptInvitation тестирует принятие приглашения
func (suite *testSuite) Test_Invitations_AcceptInvitation() {
	// Создать моки
	mockEventsPublisher := mockEvents.NewPublisher(suite.T())
	mockEventsPublisher.
		On("Publish", mock.Anything).
		Return(nil)

	usecase := &AcceptInvitationUsecase{
		Repo:            suite.RR.Chats,
		EventsPublisher: mockEventsPublisher,
	}

	suite.Run("приглашение должно существовать", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		p := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Принять приглашение
		input := In{
			SubjectID:    p.UserID,
			InvitationID: uuid.New(),
		}
		out, err := usecase.AcceptInvitation(input)
		suite.ErrorIs(err, ErrInvitationNotExists)
		suite.Zero(out)
	})

	suite.Run("приняв приглашение, пользователь становится участником чата", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		p := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(p.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Принять приглашение
		input := In{
			SubjectID:    invitation.RecipientID,
			InvitationID: invitation.ID,
		}
		out, err := usecase.AcceptInvitation(input)
		suite.Zero(out)
		suite.Require().NoError(err)
		// Получить список участников
		chats, err := suite.RR.Chats.List(chatt.Filter{})
		suite.NoError(err)
		// В списке будет 3 участника: адм., приглашаемый, приглашающий
		suite.Require().Len(chats, 1)
		suite.Require().Len(chats[0].Participants, 3)
		suite.Contains(chats[0].Participants, p)
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		var publishedEvents *events.Events

		// Создать мок
		mockEventsPublisher := mockEvents.NewPublisher(suite.T())
		mockEventsPublisher.
			On("Publish", mock.Anything).
			Run(func(args mock.Arguments) {
				publishedEvents = args.Get(0).(*events.Events)
			}).
			Return(nil)
		usecase := &AcceptInvitationUsecase{
			Repo:            suite.RR.Chats,
			EventsPublisher: mockEventsPublisher,
		}

		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		p := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(p.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Принять приглашение
		input := In{
			SubjectID:    invitation.RecipientID,
			InvitationID: invitation.ID,
		}
		out, err := usecase.AcceptInvitation(input)
		suite.Zero(out)
		suite.Require().NoError(err)

		// Проверить список опубликованных событий
		suite.Require().Len(publishedEvents.Events(), 2)

		// Событие Удаленного приглашения
		invitationRemoved := lo.FindOrElse(publishedEvents.Events(), nil, func(e any) bool {
			_, ok := e.(chatt.EventInvitationRemoved)
			return ok
		}).(chatt.EventInvitationRemoved)
		// Содержит нужных получателей
		suite.Contains(invitationRemoved.Recipients, chat.ChiefID)
		suite.Contains(invitationRemoved.Recipients, invitation.RecipientID)
		suite.Contains(invitationRemoved.Recipients, invitation.SubjectID)
		// Содержит нужное приглашение
		suite.Equal(invitation, invitationRemoved.Invitation)

		// Событие Добавленного участника
		participantAdded := lo.FindOrElse(publishedEvents.Events(), nil, func(e any) bool {
			_, ok := e.(chatt.EventParticipantAdded)
			return ok
		}).(chatt.EventParticipantAdded)
		// Содержит нужных получателей
		suite.Contains(participantAdded.Recipients, chat.ChiefID)
		suite.Contains(participantAdded.Recipients, invitation.RecipientID)
		suite.Contains(participantAdded.Recipients, invitation.SubjectID)
		// Связано с чатом
		suite.Equal(chat.ID, participantAdded.ChatID)
		// Содержит нужного участника
		suite.Equal(invitation.RecipientID, participantAdded.Participant.UserID)
	})
}
