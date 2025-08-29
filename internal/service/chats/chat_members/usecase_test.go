package chatMembers

import (
	"testing"

	"github.com/google/uuid"
	testifySuite "github.com/stretchr/testify/suite"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	serviceSuite "github.com/nice-pea/npchat/internal/service/suite"
)

type testSuite struct {
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

// Test_Members_ChatMembers тестирует получение списка участников чата
func (suite *testSuite) Test_Members_ChatMembers() {
	usecase := &ChatMembersUsecase{
		Repo: suite.RR.Chats,
	}

	suite.Run("чат должен существовать", func() {
		input := In{
			ChatID:    uuid.New(),
			SubjectID: uuid.New(),
		}
		out, err := usecase.ChatMembers(input)
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Empty(out)
	})

	suite.Run("пользователь должен быть участником чата", func() {
		// Создать чат
		chat := suite.UpsertChat(suite.RndChat())
		// Запросить список участников чата
		input := In{
			ChatID:    chat.ID,
			SubjectID: uuid.New(),
		}
		out, err := usecase.ChatMembers(input)
		// Вернется ошибка, потому пользователь не является участником чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
		suite.Empty(out)
	})

	suite.Run("возвращается список участников чата", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать несколько участников в чате
		const membersAllCount = 20
		participants := make([]chatt.Participant, membersAllCount-1)
		for i := range participants {
			// Создать участника в чате
			participants[i] = suite.AddRndParticipant(&chat)
		}
		// Сохранить чат
		suite.UpsertChat(chat)
		// Запрашивать список будет первый участник
		participant := participants[0]
		// Получить список участников в чате
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
		}
		out, err := usecase.ChatMembers(input)
		suite.NoError(err)
		suite.Require().Len(out.Participants, membersAllCount)
		// Сравнить каждого сохраненного участника с ранее созданным
		for i := range participants {
			suite.Contains(out.Participants, participants[i])
		}
	})
}
