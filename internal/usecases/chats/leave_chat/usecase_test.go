package leaveChat

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

// Test_Members_LeaveChat тестирует выход участника из чата
func (suite *testSuite) Test_Members_LeaveChat() {
	usecase := &LeaveChatUsecase{
		Repo:            suite.RR.Chats,
		EventsPublisher: mockEvents.NewPublisher(suite.T()),
	}
	// Настройка мока
	usecase.EventsPublisher.(*mockEvents.Publisher).
		On("Publish", mock.Anything).
		Return(nil)

	suite.Run("чат должен существовать", func() {
		// Покинуть чат
		input := In{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
		}
		out, err := usecase.LeaveChat(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Zero(out)
	})

	suite.Run("пользователь должен быть участником чата", func() {
		// Создать чат
		chat := suite.UpsertChat(suite.RndChat())
		// Покинуть чат
		input := In{
			SubjectID: uuid.New(),
			ChatID:    chat.ID,
		}
		out, err := usecase.LeaveChat(input)
		// Вернется ошибка, потому что пользователь не участник чата
		suite.ErrorIs(err, chatt.ErrParticipantNotExists)
		suite.Zero(out)
	})

	suite.Run("пользователь не должен быть главным администратором чата", func() {
		// Создать чат
		chat := suite.UpsertChat(suite.RndChat())
		// Покинуть чат
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
		}
		out, err := usecase.LeaveChat(input)
		// Вернется ошибка, потому что пользователь главный администратор чата
		suite.ErrorIs(err, chatt.ErrCannotRemoveChief)
		suite.Zero(out)
	})

	suite.Run("после выхода пользователь перестает быть участником", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника в этом чате
		participant := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Покинуть чат
		input := In{
			SubjectID: participant.UserID,
			ChatID:    chat.ID,
		}
		out, err := usecase.LeaveChat(input)
		suite.Require().NoError(err)
		suite.Zero(out)
		// Получить список участников чата
		filter := chatt.Filter{ParticipantID: participant.UserID}
		chats, err := suite.RR.Chats.List(filter)
		suite.Require().NoError(err)
		suite.Zero(chats)
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Новый экземпляр usecase
		usecase := &LeaveChatUsecase{
			Repo:            suite.RR.Chats,
			EventsPublisher: mockEvents.NewPublisher(suite.T()),
		}

		// Настройка мока
		var publishedEvents *events.Events
		usecase.EventsPublisher.(*mockEvents.Publisher).
			On("Publish", mock.Anything).
			Run(func(args mock.Arguments) {
				publishedEvents = args.Get(0).(*events.Events)
			}).
			Return(nil)

		// Создать чат
		chat := suite.RndChat()
		// Создать участника в этом чате
		participant := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Покинуть чат
		input := In{
			SubjectID: participant.UserID,
			ChatID:    chat.ID,
		}
		out, err := usecase.LeaveChat(input)
		suite.Require().NoError(err)
		suite.Zero(out)

		// Событие Удаленного
		participantRemoved := publishedEvents.Events()[0].(chatt.EventParticipantRemoved)
		// Содержит нужных получателей
		suite.Contains(participantRemoved.Recipients, chat.ChiefID)
		suite.Contains(participantRemoved.Recipients, participant.UserID)
		// Связано с чатом
		suite.Equal(chat.ID, participantRemoved.ChatID)
		// Содержит нужного участника
		suite.Equal(participant, participantRemoved.Participant)
	})

	suite.Run("отправленные приглашения участника, отменятся вместе с его выходом", func() {
		// TODO
	})
}
