package service

import (
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain/chatt"
)

// Test_Members_ChatMembers тестирует получение списка участников чата
func (suite *testSuite) Test_Members_ChatMembers() {
	suite.Run("чат должен существовать", func() {
		input := ChatMembersIn{
			ChatID:    uuid.New(),
			SubjectID: uuid.New(),
		}
		out, err := suite.ss.chats.ChatMembers(input)
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Empty(out)
	})

	suite.Run("пользователь должен быть участником чата", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Запросить список участников чата
		input := ChatMembersIn{
			ChatID:    chat.ID,
			SubjectID: uuid.New(),
		}
		out, err := suite.ss.chats.ChatMembers(input)
		// Вернется ошибка, потому пользователь не является участником чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
		suite.Empty(out)
	})

	suite.Run("возвращается список участников чата", func() {
		// Создать чат
		chat := suite.rndChat()
		// Создать несколько участников в чате
		const membersAllCount = 20
		participants := make([]chatt.Participant, membersAllCount-1)
		for i := range participants {
			// Создать участника в чате
			participants[i] = suite.addRndParticipant(&chat)
		}
		// Сохранить чат
		suite.upsertChat(chat)
		// Запрашивать список будет первый участник
		participant := participants[0]
		// Получить список участников в чате
		input := ChatMembersIn{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
		}
		out, err := suite.ss.chats.ChatMembers(input)
		suite.NoError(err)
		suite.Require().Len(out.Participants, membersAllCount)
		// Сравнить каждого сохраненного участника с ранее созданным
		for i := range participants {
			suite.Contains(out.Participants, participants[i])
		}
	})
}

// Test_Members_LeaveChat тестирует выход участника из чата
func (suite *testSuite) Test_Members_LeaveChat() {
	suite.Run("чат должен существовать", func() {
		// Покинуть чат
		input := LeaveChatIn{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
		}
		err := suite.ss.chats.LeaveChat(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
	})

	suite.Run("пользователь должен быть участником чата", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Покинуть чат
		input := LeaveChatIn{
			SubjectID: uuid.New(),
			ChatID:    chat.ID,
		}
		err := suite.ss.chats.LeaveChat(input)
		// Вернется ошибка, потому что пользователь не участник чата
		suite.ErrorIs(err, chatt.ErrParticipantNotExists)
	})

	suite.Run("пользователь не должен быть главным администратором чата", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Покинуть чат
		input := LeaveChatIn{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
		}
		err := suite.ss.chats.LeaveChat(input)
		// Вернется ошибка, потому что пользователь главный администратор чата
		suite.ErrorIs(err, chatt.ErrCannotRemoveChief)
	})

	suite.Run("после выхода пользователь перестает быть участником", func() {
		// Создать чат
		chat := suite.rndChat()
		// Создать участника в этом чате
		participant := suite.addRndParticipant(&chat)
		// Сохранить чат
		suite.upsertChat(chat)
		// Покинуть чат
		input := LeaveChatIn{
			SubjectID: participant.UserID,
			ChatID:    chat.ID,
		}
		err := suite.ss.chats.LeaveChat(input)
		suite.Require().NoError(err)
		// Получить список участников чата
		filter := chatt.Filter{ParticipantID: participant.UserID}
		chats, err := suite.rr.chats.List(filter)
		suite.Require().NoError(err)
		suite.Zero(chats)
	})
}

// Test_Members_DeleteMember тестирует удаление участника чата
func (suite *testSuite) Test_Members_DeleteMember() {
	suite.Run("нельзя удалить самого себя", func() {
		// Удалить участника
		userID := uuid.New()
		input := DeleteMemberIn{
			SubjectID: userID,
			ChatID:    uuid.New(),
			UserID:    userID,
		}
		err := suite.ss.chats.DeleteMember(input)
		// Вернется ошибка, потому что пользователь пытается удалить самого себя
		suite.ErrorIs(err, ErrMemberCannotDeleteHimself)
	})

	suite.Run("чат должен существовать", func() {
		// Удалить участника
		input := DeleteMemberIn{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
			UserID:    uuid.New(),
		}
		err := suite.ss.chats.DeleteMember(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
	})

	//suite.Run("subject должен быть участником чата", func() {
	//	// Создать чат
	//	chat := suite.upsertChat(suite.rndChat())
	//	// Удалить участника
	//	input := DeleteMemberIn{
	//		SubjectID: uuid.New(),
	//		ChatID:    chat.ID,
	//		UserID:    uuid.New(),
	//	}
	//	err := suite.ss.chats.DeleteMember(input)
	//	// Вернется ошибка, потому что пользователь не участник чата
	//	suite.ErrorIs(err, ErrSubjectIsNotMember)
	//})

	suite.Run("subject должен быть главным администратором чата", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		participant := suite.addRndParticipant(&chat)
		// Удалить участника
		input := DeleteMemberIn{
			SubjectID: participant.UserID,
			ChatID:    chat.ID,
			UserID:    uuid.New(),
		}
		err := suite.ss.chats.DeleteMember(input)
		// Вернется ошибка, потому что участник не главный администратор
		suite.ErrorIs(err, ErrSubjectUserIsNotChief)
	})

	suite.Run("user должен быть участником чата", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Удалить участника
		input := DeleteMemberIn{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			UserID:    uuid.New(),
		}
		err := suite.ss.chats.DeleteMember(input)
		// Вернется ошибка, потому что удаляемый пользователь не является участником
		suite.ErrorIs(err, chatt.ErrParticipantNotExists)
	})

	suite.Run("после удаления участник перестает быть участником", func() {
		// Создать чат
		chat := suite.rndChat()
		// Создать участника для удаления
		participant := suite.addRndParticipant(&chat)
		// Сохранить чат
		suite.upsertChat(chat)
		// Удалить участника
		input := DeleteMemberIn{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
			UserID:    participant.UserID,
		}
		err := suite.ss.chats.DeleteMember(input)
		suite.Require().NoError(err)
		// Найти удаленного участника
		filter := chatt.Filter{ParticipantID: participant.UserID}
		chats, err := suite.rr.chats.List(filter)
		suite.Require().NoError(err)
		// Чатов с таким пользователем нет
		suite.Empty(chats)
	})
}
