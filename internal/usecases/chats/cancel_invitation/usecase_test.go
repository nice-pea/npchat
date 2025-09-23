package cancelInvitation

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

// Test_Invitations_CancelInvitation тестирует отмену приглашения
func (suite *testSuite) Test_Invitations_CancelInvitation() {
	usecase := &CancelInvitationUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	// Настройка мока
	usecase.EventConsumer.(*mockEvents.Consumer).
		On("Consume", mock.Anything).
		Return().
		Maybe()

	suite.Run("приглашение должно существовать", func() {
		// Отменить приглашение
		input := In{
			SubjectID:    uuid.New(),
			InvitationID: uuid.New(),
		}
		out, err := usecase.CancelInvitation(input)
		// Вернется ошибка, потому что приглашения не существует
		suite.ErrorIs(err, ErrInvitationNotExists)
		suite.Zero(out)
	})

	suite.Run("приглашение могут отменить только пригласивший и приглашаемый пользователи, и администратор чата", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
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
			chat, err := chatt.Find(suite.RR.Chats, chatt.Filter{})
			suite.Require().NoError(err)
			// Создать приглашение
			invitation := suite.NewInvitation(participant.UserID, recipientID)
			suite.AddInvitation(&chat, invitation)
			// Сохранить чат
			suite.UpsertChat(chat)
			// Отменить приглашение
			input := In{
				SubjectID:    subjectUserID,
				InvitationID: invitation.ID,
			}
			out, err := usecase.CancelInvitation(input)
			suite.NoError(err)
			suite.Zero(out)
		}
	})

	suite.Run("другие участники не могут отменять приглашать ", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Случайный участник
		participantOther := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(participant.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отменить приглашение
		input := In{
			SubjectID:    participantOther.UserID,
			InvitationID: invitation.ID,
		}
		out, err := usecase.CancelInvitation(input)
		// Вернется ошибка, потому что случайный участник не может отменять приглашение
		suite.ErrorIs(err, ErrSubjectUserNotAllowed)
		suite.Zero(out)
	})

	suite.Run("после отмены, приглашение удаляется", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(participant.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отменить приглашение
		input := In{
			SubjectID:    invitation.SubjectID,
			InvitationID: invitation.ID,
		}
		out, err := usecase.CancelInvitation(input)
		suite.Require().NoError(err)
		suite.Zero(out)
		// Получить список приглашений
		chats, err := suite.RR.Chats.List(chatt.Filter{})
		suite.NoError(err)
		suite.Require().Len(chats, 1)
		suite.Empty(chats[0].Invitations)
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Новый экземпляр usecase
		usecase := &CancelInvitationUsecase{
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
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(participant.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отменить приглашение
		input := In{
			SubjectID:    invitation.SubjectID,
			InvitationID: invitation.ID,
		}
		out, err := usecase.CancelInvitation(input)
		suite.Require().NoError(err)
		suite.Zero(out)

		// Проверить список опубликованных событий
		suite.AssertHasEventType(consumedEvents, chatt.EventInvitationRemoved)
	})
}
