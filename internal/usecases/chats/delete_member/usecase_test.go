package deleteMember

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

// Test_Members_DeleteMember тестирует удаление участника чата
func (suite *testSuite) Test_Members_DeleteMember() {
	usecase := &DeleteMemberUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	// Настройка мока
	usecase.EventConsumer.(*mockEvents.Consumer).
		On("Consume", mock.Anything).
		Return()

	suite.Run("нельзя удалить самого себя", func() {
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
		// Удалить участника
		input := In{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
			UserID:    uuid.New(),
		}
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
		// Создать чат
		chat := suite.UpsertChat(suite.RndChat())
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Удалить участника
		input := In{
			SubjectID: participant.UserID,
			ChatID:    chat.ID,
			UserID:    uuid.New(),
		}
		out, err := usecase.DeleteMember(input)
		// Вернется ошибка, потому что участник не главный администратор
		suite.ErrorIs(err, ErrSubjectUserIsNotChief)
		suite.Zero(out)
	})

	suite.Run("user должен быть участником чата", func() {
		// Создать чат
		chat := suite.UpsertChat(suite.RndChat())
		// Удалить участника
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			UserID:    uuid.New(),
		}
		out, err := usecase.DeleteMember(input)
		// Вернется ошибка, потому что удаляемый пользователь не является участником
		suite.ErrorIs(err, chatt.ErrParticipantNotExists)
		suite.Zero(out)
	})

	suite.Run("после удаления участник перестает быть участником", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника для удаления
		participant := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Удалить участника
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			UserID:    participant.UserID,
		}
		out, err := usecase.DeleteMember(input)
		suite.Require().NoError(err)
		suite.Zero(out)
		// Найти удаленного участника
		filter := chatt.Filter{ParticipantID: participant.UserID}
		chats, err := suite.RR.Chats.List(filter)
		suite.Require().NoError(err)
		// Чатов с таким пользователем нет
		suite.Empty(chats)
	})

	suite.Run("после завершения операции, будут созданы события", func() {
		// Новый экземпляр usecase
		usecase := &DeleteMemberUsecase{
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
		// Создать участника для удаления
		participant := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Удалить участника
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			UserID:    participant.UserID,
		}
		out, err := usecase.DeleteMember(input)
		suite.Require().NoError(err)
		suite.Zero(out)

		// Проверить список опубликованных событий
		suite.True(serviceSuite.HasElementOfType[chatt.EventParticipantRemoved](consumedEvents))
	})
}
