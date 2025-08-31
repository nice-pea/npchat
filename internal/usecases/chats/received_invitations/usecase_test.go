package receivedInvitations

import (
	"testing"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	serviceSuite "github.com/nice-pea/npchat/internal/usecases/suite"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

// Test_Invitations_ReceivedInvitations тестирует получение списка приглашений направленных пользователю
func (suite *testSuite) Test_Invitations_ReceivedInvitations() {
	usecase := &ReceivedInvitationsUsecase{
		Repo: suite.RR.Chats,
	}

	suite.Run("пользователя не приглашали и потому вернется пустой список", func() {
		// Получить список приглашений
		input := In{
			SubjectID: uuid.New(),
		}
		invitations, err := usecase.ReceivedInvitations(input)
		suite.NoError(err)
		suite.Empty(invitations)
	})

	suite.Run("из многих приглашений несколько направлено пользователю и потому вернутся только они", func() {
		// ID пользователя
		userID := uuid.New()
		// Создать несколько приглашений, направленных пользователю
		invitationsOfUser := make([]chatt.Invitation, 5)
		for i := range invitationsOfUser {
			// Создать чат
			chat := suite.RndChat()
			p := suite.AddRndParticipant(&chat)
			invitationsOfUser[i] = suite.NewInvitation(p.UserID, userID)
			suite.AddInvitation(&chat, invitationsOfUser[i])
			// Сохранить чат
			suite.UpsertChat(chat)
		}
		// Создать несколько приглашений, направленных каким-то другим пользователям
		for range 10 {
			// Создать чат
			chat := suite.RndChat()
			p := suite.AddRndParticipant(&chat)
			i := suite.NewInvitation(p.UserID, uuid.New())
			suite.AddInvitation(&chat, i)
			// Сохранить чат
			suite.UpsertChat(chat)
		}
		// Получить список приглашений
		input := In{
			SubjectID: userID,
		}
		out, err := usecase.ReceivedInvitations(input)
		suite.NoError(err)
		// В списке будут только приглашения, направленные пользователю
		suite.Len(out.ChatsInvitations, len(invitationsOfUser))
		for _, invitation := range out.ChatsInvitations {
			suite.Contains(invitationsOfUser, invitation)
		}
	})
}
