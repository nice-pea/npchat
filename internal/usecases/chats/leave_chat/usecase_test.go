package leaveChat

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/chatt"
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
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	// Настройка мока
	usecase.EventConsumer.(*mockEvents.Consumer).
		On("Consume", mock.Anything).
		Return()

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
			Repo:          suite.RR.Chats,
			EventConsumer: mockEvents.NewConsumer(suite.T()),
		}
		// Настройка мока
		var consumedEvents []any
		usecase.EventConsumer.(*mockEvents.Consumer).
			On("Consume", mock.Anything).
			Run(func(args mock.Arguments) {
				consumedEvents = append(consumedEvents, args.Get(0).([]any)...)
			}).
			Return()

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

		// Проверить список опубликованных событий
		suite.True(serviceSuite.HasElementOfType[chatt.EventParticipantRemoved](consumedEvents))
	})

	suite.Run("отправленные приглашения участника, отменятся вместе с его выходом", func() {
		// TODO
	})
}
