package chatMembers

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
	serviceSuite.Suite
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(testSuite))
}

// Test_Members_ChatMembers тестирует получение списка участников чата
func (suite *testSuite) Test_Members_ChatMembers() {
	suite.Run("чат должен существовать", func() {
		// Создать usecase и моки
		usecase, mockRepo := newUsecase(suite)
		input := In{
			ChatID:    uuid.New(),
			SubjectID: uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{}, nil).Once()
		out, err := usecase.ChatMembers(input)
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Empty(out)
	})

	suite.Run("пользователь должен быть участником чата", func() {
		// Создать usecase и моки
		usecase, mockRepo := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Запросить список участников чата
		input := In{
			ChatID:    chat.ID,
			SubjectID: uuid.New(),
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()
		out, err := usecase.ChatMembers(input)
		// Вернется ошибка, потому пользователь не является участником чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
		suite.Empty(out)
	})

	suite.Run("возвращается список участников чата", func() {
		// Создать usecase и моки
		usecase, mockRepo := newUsecase(suite)
		// Создать чат
		chat := suite.RndChat()
		// Создать несколько участников в чате
		const membersAllCount = 20
		participants := make([]chatt.Participant, membersAllCount-1)
		for i := range participants {
			// Создать участника в чате
			participants[i] = suite.AddRndParticipant(&chat)
		}
		// Запрашивать список будет первый участник
		participant := participants[0]
		// Получить список участников в чате
		input := In{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
		}
		mockRepo.EXPECT().List(mock.Anything).Return([]chatt.Chat{chat}, nil).Once()

		out, err := usecase.ChatMembers(input)
		suite.NoError(err)
		suite.Require().Len(out.Participants, membersAllCount)
		// Сравнить каждого сохраненного участника с ранее созданным
		for i := range participants {
			suite.Contains(out.Participants, participants[i])
		}
	})
}

func newUsecase(suite *testSuite) (*ChatMembersUsecase, *mockChatt.Repository) {
	uc := &ChatMembersUsecase{
		Repo: suite.RR.Chats,
	}
	mockRepo := uc.Repo.(*mockChatt.Repository)
	return uc, mockRepo
}
