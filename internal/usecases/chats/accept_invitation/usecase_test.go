package acceptInvitation

import (
	"slices"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func Test_Invitations_AcceptInvitation_Testing(t *testing.T) {
	newUsecase := func() (*AcceptInvitationUsecase, *mockChatt.Repository, *mockEvents.Consumer) {
		mockRepo := mockChatt.NewRepository(t)
		mockEvent := mockEvents.NewConsumer(t)
		usecase := &AcceptInvitationUsecase{
			Repo:          mockRepo,
			EventConsumer: mockEvent,
		}
		// Настройка мока

		return usecase, mockRepo, mockEvent
	}

	t.Run("приглашение должно существовать", func(t *testing.T) {
		usecase, mockRepo, mockEvent := newUsecase()
		mockEvent.EXPECT().Consume(mock.Anything).Return().Maybe()
		// Создать чат
		chat := RndChat(t)
		// Создать участника
		p := AddRndParticipant(t, &chat)
		// Сохранить чат
		mockRepo.EXPECT().Upsert(chat).Return(nil)
		UpsertChat(t, mockRepo, chat)
		// Принять приглашение
		input := In{
			SubjectID:    p.UserID,
			InvitationID: uuid.New(),
		}
		mockRepo.EXPECT().List(chatt.Filter{
			InvitationID: input.InvitationID,
		}).Return(nil, chatt.ErrChatNotExists)

		out, err := usecase.AcceptInvitation(input)
		assert.ErrorIs(t, err, ErrInvitationNotExists)
		assert.Zero(t, out)
	})
	t.Run("приняв приглашение, пользователь становится участником чата", func(t *testing.T) {
		usecase, mockRepo, mockEvent := newUsecase()
		mockEvent.EXPECT().Consume(mock.Anything).Return().Maybe()
		// Создать чат
		chat := RndChat(t)
		// Создать участника
		p := AddRndParticipant(t, &chat)
		// Создать приглашение
		invitation := NewInvitation(t, p.UserID, uuid.New())
		AddInvitation(t, &chat, invitation)
		// Сохранить чат
		mockRepo.EXPECT().Upsert(chat).Return(nil)
		UpsertChat(t, mockRepo, chat)
		// Принять приглашение
		input := In{
			SubjectID:    invitation.RecipientID,
			InvitationID: invitation.ID,
		}
		mockRepo.EXPECT().List(chatt.Filter{
			InvitationID: input.InvitationID,
		}).Return([]chatt.Chat{chat}, nil).Once()
		mockRepo.EXPECT().Upsert(mock.Anything).RunAndReturn(func(updatedChat chatt.Chat) error {
			// Базовые проверки
			assert.Equal(t, chat.ID, updatedChat.ID)
			// Дополнительные проверки логики можно добавить здесь
			return nil
		}).Once()

		out, err := usecase.AcceptInvitation(input)
		assert.Zero(t, out)
		require.NoError(t, err)
		// Получить список участников
		mockRepo.EXPECT().List(chatt.Filter{}).RunAndReturn(func(filter chatt.Filter) ([]chatt.Chat, error) {
			chat.Participants = append(chat.Participants, chatt.Participant{
				UserID: input.InvitationID,
			})
			return []chatt.Chat{chat}, nil
		}).Once()
		chats, err := mockRepo.List(chatt.Filter{})
		assert.NoError(t, err)
		// В списке будет 3 участника: адм., приглашаемый, приглашающий
		require.Len(t, chats, 1)
		require.Len(t, chats[0].Participants, 3)
		require.Contains(t, chats[0].Participants, p)
	})
	t.Run("после завершения операции, будут созданы события", func(t *testing.T) {
		// Новый экземпляр usecase
		usecase, mockRepo, mockEvent := newUsecase()
		// Настройка мока
		var consumedEvents []events.Event
		mockEvent.EXPECT().Consume(mock.Anything).Run(
			func(args []events.Event) {
				consumedEvents = append(consumedEvents, args...)
			}).Return()
		// Создать чат
		chat := RndChat(t)
		// Создать участника
		p := AddRndParticipant(t, &chat)
		// Создать приглашение
		invitation := NewInvitation(t, p.UserID, uuid.New())
		AddInvitation(t, &chat, invitation)
		// Сохранить чат
		mockRepo.EXPECT().Upsert(chat).Return(nil)
		UpsertChat(t, mockRepo, chat)
		// Принять приглашение
		input := In{
			SubjectID:    invitation.RecipientID,
			InvitationID: invitation.ID,
		}
		mockRepo.EXPECT().List(chatt.Filter{
			InvitationID: input.InvitationID,
		}).Return([]chatt.Chat{chat}, nil).Once()
		mockRepo.EXPECT().Upsert(mock.Anything).RunAndReturn(func(updatedChat chatt.Chat) error {
			// Базовые проверки
			assert.Equal(t, chat.ID, updatedChat.ID)
			// Дополнительные проверки логики можно добавить здесь
			return nil
		}).Once()

		out, err := usecase.AcceptInvitation(input)
		assert.Zero(t, out)
		require.NoError(t, err)
		// Проверить список опубликованных событий
		AssertHasEventType(t, consumedEvents, chatt.EventInvitationRemoved)
		AssertHasEventType(t, consumedEvents, chatt.EventParticipantAdded)
	})
}

func AssertHasEventType(t *testing.T, ee []events.Event, eventType string) {
	t.Helper()
	assert.True(t, slices.ContainsFunc(ee, func(e events.Event) bool {
		return e.Type == eventType
	}))
}

func AddInvitation(t *testing.T, chat *chatt.Chat, invitation chatt.Invitation) {
	require.NoError(t, chat.AddInvitation(invitation, nil))
}

func NewInvitation(t *testing.T, subjectID, recipientID uuid.UUID) chatt.Invitation {
	i, err := chatt.NewInvitation(subjectID, recipientID)
	require.NoError(t, err)
	return i
}

func RndChat(t *testing.T) chatt.Chat {
	chat, err := chatt.NewChat(gofakeit.Noun(), uuid.New(), nil)
	require.NoError(t, err)
	return chat
}

func AddRndParticipant(t *testing.T, chat *chatt.Chat) chatt.Participant {
	p, err := chatt.NewParticipant(uuid.New())
	require.NoError(t, err)
	require.NoError(t, chat.AddParticipant(p, nil))

	return p
}

func UpsertChat(t *testing.T, mockRepo *mockChatt.Repository, chat chatt.Chat) chatt.Chat {
	err := mockRepo.Upsert(chat)
	require.NoError(t, err)

	return chat
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

// Test_Invitations_AcceptInvitation тестирует принятие приглашения
func (suite *testSuite) Test_Invitations_AcceptInvitation() {
	usecase := &AcceptInvitationUsecase{
		Repo:          suite.RR.Chats,
		EventConsumer: mockEvents.NewConsumer(suite.T()),
	}
	// Настройка мока
	usecase.EventConsumer.(*mockEvents.Consumer).
		On("Consume", mock.Anything).
		Return().
		Maybe()

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
		// Новый экземпляр usecase
		usecase := &AcceptInvitationUsecase{
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
		suite.AssertHasEventType(consumedEvents, chatt.EventInvitationRemoved)
		suite.AssertHasEventType(consumedEvents, chatt.EventParticipantAdded)
	})
}
