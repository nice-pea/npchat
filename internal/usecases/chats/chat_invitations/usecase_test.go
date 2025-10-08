package chatInvitations

import (
	"testing"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	mockChatt "github.com/nice-pea/npchat/internal/domain/chatt/mocks"
	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
)

type testSuite struct {
	serviceSuite.SuiteWithMocks
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

// Test_Invitations_ChatInvitations тестирует получение списка приглашений
func (suite *testSuite) Test_Invitations_ChatInvitations() {

	suite.Run("чат должен существовать", func() {
		// Создать usecase и моки
		usecase, mockRepo := newUsecase(suite)
		// Получить список приглашений
		input := In{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
		}
		mockRepo.EXPECT().List(chatt.Filter{ID: input.ChatID}).Return([]chatt.Chat{}, nil)
		out, err := usecase.ChatInvitations(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Empty(out.Invitations)
	})

	suite.Run("субъект должен быть участником чата", func() {
		// Создать usecase и моки
		usecase, mockRepo := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Получить список приглашений
		input := In{
			ChatID:    chat.ID,
			SubjectID: uuid.New(),
		}
		mockRepo.EXPECT().List(chatt.Filter{ID: input.ChatID}).Return([]chatt.Chat{chat}, nil)
		out, err := usecase.ChatInvitations(input)
		// Вернется ошибка, потому что пользователь не участник чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
		suite.Empty(out.Invitations)
	})

	suite.Run("пустой список из чата без приглашений", func() {
		// Создать usecase и моки
		usecase, mockRepo := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Получить список приглашений
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
		}
		mockRepo.EXPECT().List(chatt.Filter{ID: input.ChatID}).Return([]chatt.Chat{chat}, nil)
		out, err := usecase.ChatInvitations(input)
		suite.NoError(err)
		suite.Empty(out.Invitations)
	})

	suite.Run("субъект не администратор чата и видит только отправленные им приглашения", func() {
		// Создать usecase и моки
		usecase, mockRepo := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		participant := suite.AddRndParticipant(&chat)
		// Создать приглашения, отправленные участником
		subjectInvitations := make([]chatt.Invitation, 3)
		for i := range subjectInvitations {
			subjectInvitations[i] = suite.NewInvitation(participant.UserID, uuid.New())
			suite.AddInvitation(&chat, subjectInvitations[i])
		}
		// Создать приглашения, отправленные какими-то другими пользователями
		for range 3 {
			p := suite.AddRndParticipant(&chat)
			i := suite.NewInvitation(p.UserID, uuid.New())
			suite.AddInvitation(&chat, i)
		}
		// Получить список приглашений
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
		}
		mockRepo.EXPECT().List(chatt.Filter{ID: input.ChatID}).Return([]chatt.Chat{chat}, nil)
		out, err := usecase.ChatInvitations(input)
		suite.Require().NoError(err)
		// В списке будут приглашения, отправленные участником
		suite.Len(out.Invitations, len(subjectInvitations))
		for i, subjectInvitation := range subjectInvitations {
			suite.Equal(subjectInvitation, out.Invitations[i])
		}
	})

	suite.Run("субъект является администратором чата и видит все отправленные приглашения в чат", func() {
		// Создать usecase и моки
		usecase, mockRepo := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Создать приглашения, отправленные какими-то пользователями
		invitationsSent := make([]chatt.Invitation, 3)
		for i := range invitationsSent {
			p := suite.AddRndParticipant(&chat)
			invitationsSent[i] = suite.NewInvitation(p.UserID, uuid.New())
			suite.AddInvitation(&chat, invitationsSent[i])
		}
		// Получить список приглашений
		input := In{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
		}
		mockRepo.EXPECT().List(chatt.Filter{ID: input.ChatID}).Return([]chatt.Chat{chat}, nil)
		out, err := usecase.ChatInvitations(input)
		suite.Require().NoError(err)
		// В списке будут приглашения все приглашения
		suite.Len(out.Invitations, len(invitationsSent))
		for _, saved := range invitationsSent {
			suite.Contains(out.Invitations, saved)
		}
	})
}

func newUsecase(suite *testSuite) (*ChatInvitationsUsecase, *mockChatt.Repository) {
	uc := &ChatInvitationsUsecase{
		Repo: suite.RR.Chats,
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	return uc, mockRepo
}
