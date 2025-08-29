package service

import (
	"github.com/google/uuid"
)

func (suite *testSuite) newUserChatsInput(userID uuid.UUID) WhichParticipateIn {
	return WhichParticipateIn{
		SubjectID: userID,
		UserID:    userID,
	}
}

// Test_Chats_UserChats тестирует запрос список чатов в которых участвует пользователь
func (suite *testSuite) Test_Chats_UserChats() {
	suite.Run("пользователь может запрашивать только свой чат", func() {
		input := WhichParticipateIn{
			SubjectID: uuid.New(),
			UserID:    uuid.New(),
		}
		out, err := suite.ss.chats.WhichParticipate(input)
		suite.ErrorIs(err, ErrUnauthorizedChatsView)
		suite.Empty(out)
	})

	suite.Run("пустой список из пустого репозитория", func() {
		input := suite.newUserChatsInput(uuid.New())
		out, err := suite.ss.chats.WhichParticipate(input)
		suite.NoError(err)
		suite.Empty(out)
	})

	suite.Run("пустой список если у пользователя нет чатов", func() {
		const chatsAllCount = 11
		for range chatsAllCount {
			chat := suite.UpsertChat(suite.RndChat())
			suite.AddRndParticipant(&chat)
		}
		input := suite.newUserChatsInput(uuid.New())
		out, err := suite.ss.chats.WhichParticipate(input)
		suite.NoError(err)
		suite.Empty(out)
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
		out, err := suite.ss.chats.WhichParticipate(input)
		suite.NoError(err)
		suite.Len(out.Chats, chatsAllCount)
	})
}
