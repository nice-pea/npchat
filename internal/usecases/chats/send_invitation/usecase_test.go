package sendInvitation

import (
	"testing"

	"github.com/google/uuid"
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

// Test_Invitations_SendChatInvitation тестирует отправку приглашения
func (suite *testSuite) Test_Invitations_SendChatInvitation() {
	usecase := &SendInvitationUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	// Настройка мока
	usecase.EventConsumer.(*mockEvents.Consumer).
		On("Consume", mock.Anything).
		Return().
		Maybe()

	suite.Run("чат должен существовать", func() {
		// Отправить приглашение
		input := In{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
			UserID:    uuid.New(),
		}
		invitation, err := usecase.SendInvitation(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Zero(invitation)
	})

	suite.Run("субъект должен быть участником", func() {
		// Создать чат
		chat := suite.RndChat()
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отправить приглашение
		input := In{
			SubjectID: uuid.New(),
			ChatID:    chat.ID,
			UserID:    uuid.New(),
		}
		invitation, err := usecase.SendInvitation(input)
		// Вернется ошибка, потому что субъект не является участником чата
		suite.ErrorIs(err, chatt.ErrSubjectIsNotMember)
		suite.Zero(invitation)
	})

	suite.Run("приглашаемый пользователь может не существовать", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отправить приглашение
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    uuid.New(),
		}
		out, err := usecase.SendInvitation(input)
		suite.NoError(err)
		suite.NotZero(out)
	})

	suite.Run("приглашаемый пользователь не должен состоять в этом чате", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать участника для приглашаемого пользователя
		participantInvitating := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отправить приглашение
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    participantInvitating.UserID,
		}
		invitation, err := usecase.SendInvitation(input)
		// Вернется ошибка, потому что приглашаемый пользователь уже является участником этого чата
		suite.ErrorIs(err, chatt.ErrParticipantExists)
		suite.Zero(invitation)
	})

	suite.Run("одновременно не может существовать несколько приглашений одного пользователя в этот чат", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать приглашаемого пользователя
		targetUserID := uuid.New()
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отправить приглашение
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    targetUserID,
		}
		invitation, err := usecase.SendInvitation(input)
		suite.NoError(err)
		suite.Require().NotZero(invitation)
		// Отправить повторно приглашение
		invitation, err = usecase.SendInvitation(input)
		// Вернется ошибка, потому что этот пользователь уже приглашен в чат
		suite.ErrorIs(err, chatt.ErrUserIsAlreadyInvited)
		suite.Zero(invitation)
	})

	suite.Run("любой участник может приглашать много пользователей", func() {
		// Создать чат
		chat := suite.RndChat()
		// Сохранить чат
		suite.UpsertChat(chat)
		// Создать много приглашений от разных участников
		var invitationsCreated []chatt.Invitation
		for range 5 {
			chat, err := chatt.Find(suite.RR.Chats, chatt.Filter{})
			suite.Require().NoError(err)
			// Создать участника
			participant := suite.AddRndParticipant(&chat)
			// Сохранить чат
			suite.UpsertChat(chat)
			for range 5 {
				// Отправить приглашение
				input := In{
					ChatID:    chat.ID,
					SubjectID: participant.UserID,
					UserID:    uuid.New(),
				}
				out, err := usecase.SendInvitation(input)
				suite.NoError(err)
				suite.Require().NotZero(out)
				invitationsCreated = append(invitationsCreated, out.Invitation)
			}
		}
		// Получить список приглашений
		chats, err := suite.RR.Chats.List(chatt.Filter{})
		suite.NoError(err)
		// В списке содержатся все созданные приглашения
		suite.Require().Len(chats, 1)
		suite.Require().Len(chats[0].Invitations, len(invitationsCreated))
		for _, createdInvitation := range invitationsCreated {
			suite.Contains(chats[0].Invitations, createdInvitation)
		}
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Новый экземпляр usecase
		usecase := &SendInvitationUsecase{
			Repo:          suite.RR.Chats,
			EventConsumer: mockEvents.NewConsumer(suite.T()),
		}
		// Настройка мока
		var consumedEvents []events.Event
		usecase.EventConsumer.(*mockEvents.Consumer).
			On("Consume", mock.Anything).
			Run(func(args mock.Arguments) {
				consumedEvents = append(consumedEvents, args.Get(0).([]events.Event)...)
			}).
			Return()

		// Создать чат
		chat := suite.RndChat()
		// Сохранить чат
		suite.UpsertChat(chat)
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отправить приглашение
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    uuid.New(),
		}
		out, err := usecase.SendInvitation(input)
		suite.NoError(err)
		suite.Require().NotZero(out)

		// Проверить список опубликованных событий
		suite.AssertHasEventType(consumedEvents, chatt.EventInvitationAddedType)
	})
}
