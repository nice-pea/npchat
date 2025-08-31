package createChat

import (
	"fmt"
	"math/rand/v2"
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

func (suite *testSuite) newCreateInputRandom() In {
	return In{
		ChiefUserID: uuid.New(),
		Name:        fmt.Sprintf("name%d", rand.Int()),
	}
}

// Test_Chats_CreateChat тестирует создание чата
func (suite *testSuite) Test_Chats_CreateChat() {
	usecase := &CreateChatUsecase{
		Repo: suite.RR.Chats,
	}

	suite.Run("выходящие совпадают с заданными", func() {
		// Создать чат
		input := suite.newCreateInputRandom()
		out, err := usecase.CreateChat(input)
		suite.NoError(err)
		// Сравнить результат с входящими значениями
		suite.Equal(input.Name, out.Chat.Name)
		suite.Equal(input.ChiefUserID, out.Chat.ChiefID)
	})

	suite.Run("можно затем прочитать из репозитория", func() {
		// Создать чат
		input := suite.newCreateInputRandom()
		out, err := usecase.CreateChat(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// Получить список чатов
		chats, err := suite.RR.Chats.List(chatt.Filter{})
		suite.Require().NoError(err)
		// В списке этот чат будет единственным
		suite.Require().Len(chats, 1)
		suite.Equal(out.Chat.Name, chats[0].Name)
		suite.Equal(out.Chat.ChiefID, chats[0].ChiefID)
	})

	suite.Run("создается участник для главного администратора", func() {
		// Создать чат
		input := suite.newCreateInputRandom()
		out, err := usecase.CreateChat(input)
		suite.Require().NoError(err)
		suite.Require().NotZero(out)
		// Получить список участников
		//members, err := usecase.MembersRepo.List(domain.MembersFilter{})
		//suite.NoError(err)
		// В списке этот участник будет единственным
		suite.Require().Len(out.Chat.Participants, 1)
		// Участником является главный администратор созданного чата
		suite.Equal(input.ChiefUserID, out.Chat.Participants[0].UserID)
	})

	suite.Run("можно создать чаты с одинаковым именем", func() {
		input := suite.newCreateInputRandom()
		// Создать несколько чатов с одинаковым именем
		const chatsAllCount = 2
		for range chatsAllCount {
			out, err := usecase.CreateChat(input)
			suite.Require().NoError(err)
			suite.Require().NotZero(out)
		}
		// Получить список чатов
		chats, err := suite.RR.Chats.List(chatt.Filter{})
		suite.NoError(err)
		// Количество чатов равно количеству созданных
		suite.Len(chats, chatsAllCount)
	})

	suite.Run("количество созданных чатов на одного пользователя не ограничено", func() {
		// Пользователь
		userID := uuid.New()
		// Создать много чатов от лица пользователя
		const chatsAllCount = 900
		for range chatsAllCount {
			out, err := usecase.CreateChat(In{
				ChiefUserID: userID,
				Name:        "name",
			})
			suite.Require().NoError(err)
			suite.Require().NotZero(out)
		}
		// Получить список чатов
		chats, err := suite.RR.Chats.List(chatt.Filter{})
		suite.NoError(err)
		// Количество чатов равно количеству созданных
		suite.Len(chats, chatsAllCount)
	})
}
