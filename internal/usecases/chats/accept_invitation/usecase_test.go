package acceptInvitation

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

// Test_Invitations_AcceptInvitation тестирует принятие приглашения
func (suite *testSuite) Test_Invitations_AcceptInvitation() {
	usecase := &AcceptInvitationUsecase{
		Repo: suite.RR.Chats,
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
}
