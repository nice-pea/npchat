package receivedInvitations

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
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

// Test_Invitations_ReceivedInvitations тестирует получение списка приглашений направленных пользователю
func (suite *testSuite) Test_Invitations_ReceivedInvitations() {

	suite.Run("пользователя не приглашали и потому вернется пустой список", func() {
		usecase, mockRepo := newUsecase(suite)
		// Получить список приглашений
		input := In{
			SubjectID: uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{}, nil).Once()
		invitations, err := usecase.ReceivedInvitations(input)
		suite.NoError(err)
		suite.Empty(invitations)
	})

	suite.Run("из многих приглашений несколько направлено пользователю и потому вернутся только они", func() {
		usecase, mockRepo := newUsecase(suite)
		// ID пользователя
		userID := uuid.New()
		// Создать несколько приглашений, направленных пользователю
		invitationsOfUser := make([]chatt.Invitation, 5)
		chats := make([]chatt.Chat, 0, len(invitationsOfUser))
		for i := range invitationsOfUser {
			// Создать чат
			chat := suite.RndChat()
			p := suite.AddRndParticipant(&chat)
			invitationsOfUser[i] = suite.NewInvitation(p.UserID, userID)
			suite.AddInvitation(&chat, invitationsOfUser[i])
			chats = append(chats, chat)
		}
		// Создать несколько приглашений, направленных каким-то другим пользователям
		for range 10 {
			// Создать чат
			chat := suite.RndChat()
			p := suite.AddRndParticipant(&chat)
			i := suite.NewInvitation(p.UserID, uuid.New())
			suite.AddInvitation(&chat, i)
		}
		// Получить список приглашений
		input := In{
			SubjectID: userID,
		}
		mockRepo.EXPECT().List(chatt.Filter{
			InvitationRecipientID: input.SubjectID,
		}).Return(chats, nil).Once()
		out, err := usecase.ReceivedInvitations(input)
		suite.NoError(err)
		// В списке будут только приглашения, направленные пользователю
		suite.Len(out.ChatsInvitations, len(invitationsOfUser))
		for _, invitation := range out.ChatsInvitations {
			suite.Contains(invitationsOfUser, invitation)
		}
	})

}

func newUsecase(suite *testSuite) (*ReceivedInvitationsUsecase, *mockChatt.Repository) {
	uc := &ReceivedInvitationsUsecase{
		Repo: suite.RR.Chats,
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	return uc, mockRepo
}
