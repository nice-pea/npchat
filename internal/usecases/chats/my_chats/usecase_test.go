package myChats

import (
	"testing"
	"time"

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

// Test_Chats_MyChats тестирует запрос список чатов в которых участвует пользователь
func (suite *testSuite) Test_Chats_MyChats() {
	usecase := &MyChatsUsecase{
		Repo: suite.RR.Chats,
	}

	suite.Run("пользователь может запрашивать только свой чат", func() {
		input := In{
			SubjectID: uuid.New(),
			UserID:    uuid.New(),
		}
		out, err := usecase.MyChats(input)
		suite.ErrorIs(err, ErrUnauthorizedChatsView)
		suite.Empty(out.Chats)
	})

	suite.Run("пустой список из пустого репозитория", func() {
		input := suite.newUserChatsInput(uuid.New())
		out, err := usecase.MyChats(input)
		suite.NoError(err)
		suite.Empty(out.Chats)
	})

	suite.Run("пустой список если у пользователя нет чатов", func() {
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.UpsertChat(suite.RndChat())
			suite.AddRndParticipant(&chat)
		}
		input := suite.newUserChatsInput(uuid.New())
		out, err := usecase.MyChats(input)
		suite.NoError(err)
		suite.Empty(out.Chats)
	})

	suite.Run("у пользователя может быть несколько чатов", func() {
		userID := uuid.New()
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.RndChat()
			p := suite.NewParticipant(userID)
			suite.AddParticipant(&chat, p)
			suite.UpsertChat(chat)
		}
		input := suite.newUserChatsInput(userID)
		out, err := usecase.MyChats(input)
		suite.NoError(err)
		suite.Len(out.Chats, chatsAllCount)
	})

	suite.Run("возвращает NextKeyset пустым если осталось чатов меньше чем page", func() {
		userID := uuid.New()
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.RndChat()
			p := suite.NewParticipant(userID)
			suite.AddParticipant(&chat, p)
			suite.UpsertChat(chat)
		}
		input := suite.newUserChatsInput(userID)
		out, err := usecase.MyChats(input)
		suite.NoError(err)
		suite.Zero(out.NextKeyset)
	})

	suite.Run("с установленным Keyset вернется только часть чатов", func() {
		userID := uuid.New()
		//const chatsAllCount = 11
		now := time.Now()
		activeBefore := now.Add(-time.Hour)
		// Список чатов, которые должны вернуться в результате
		expectedChats := make([]chatt.Chat, 10)
		for i := range expectedChats {
			chat := suite.RndChat()
			chat.LastActiveAt = activeBefore.
				Add(-time.Duration(i+1) * time.Minute).
				UTC().
				Truncate(time.Millisecond)
			suite.AddParticipant(&chat, suite.NewParticipant(userID))
			expectedChats[i] = suite.UpsertChat(chat)
		}
		// Чаты, которые не должны вернуться
		for range 11 {
			chat := suite.RndChat()
			suite.AddParticipant(&chat, suite.NewParticipant(userID))
			suite.UpsertChat(chat)
		}
		// Получить список чатов
		out, err := usecase.MyChats(In{
			SubjectID: userID,
			UserID:    userID,
			Keyset: Keyset{
				ActiveBefore: activeBefore,
			},
		})
		suite.NoError(err)
		suite.Zero(out.NextKeyset)
		suite.Require().Len(out.Chats, len(expectedChats))
		for i, chat := range out.Chats {
			suite.Equal(expectedChats[i], chat)
		}
	})
}

func (suite *testSuite) newUserChatsInput(userID uuid.UUID) In {
	return In{
		SubjectID: userID,
		UserID:    userID,
	}
}
