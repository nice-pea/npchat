package service

import (
	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain/chatt"
)

// Test_Invitations_ChatInvitations тестирует получение списка приглашений
func (suite *testSuite) Test_Invitations_ChatInvitations() {
	suite.Run("чат должен существовать", func() {
		// Получить список приглашений
		input := ChatInvitationsIn{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
		}
		out, err := suite.ss.chats.ChatInvitations(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Empty(out.Invitations)
	})

	suite.Run("субъект должен быть участником чата", func() {
		// Создать чат
		chat := suite.UpsertChat(suite.RndChat())
		// Получить список приглашений
		input := ChatInvitationsIn{
			ChatID:    chat.ID,
			SubjectID: uuid.New(),
		}
		out, err := suite.ss.chats.ChatInvitations(input)
		// Вернется ошибка, потому что пользователь не участник чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
		suite.Empty(out.Invitations)
	})

	suite.Run("пустой список из чата без приглашений", func() {
		// Создать чат
		chat := suite.RndChat()
		// Сохранить чат
		suite.UpsertChat(chat)
		// Получить список приглашений
		input := ChatInvitationsIn{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
		}
		out, err := suite.ss.chats.ChatInvitations(input)
		suite.NoError(err)
		suite.Empty(out.Invitations)
	})

	suite.Run("субъект не администратор чата и видит только отправленные им приглашения", func() {
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
		// Сохранить чат
		suite.UpsertChat(chat)
		// Получить список приглашений
		input := ChatInvitationsIn{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
		}
		out, err := suite.ss.chats.ChatInvitations(input)
		suite.Require().NoError(err)
		// В списке будут приглашения, отправленные участником
		suite.Len(out.Invitations, len(subjectInvitations))
		for i, subjectInvitation := range subjectInvitations {
			suite.Equal(subjectInvitation, out.Invitations[i])
		}
	})

	suite.Run("субъект является администратором чата и видит все отправленные приглашения в чат", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать приглашения, отправленные какими-то пользователями
		invitationsSent := make([]chatt.Invitation, 3)
		for i := range invitationsSent {
			p := suite.AddRndParticipant(&chat)
			invitationsSent[i] = suite.NewInvitation(p.UserID, uuid.New())
			suite.AddInvitation(&chat, invitationsSent[i])
		}
		// Сохранить чат
		suite.UpsertChat(chat)
		// Получить список приглашений
		input := ChatInvitationsIn{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
		}
		out, err := suite.ss.chats.ChatInvitations(input)
		suite.Require().NoError(err)
		// В списке будут приглашения все приглашения
		suite.Len(out.Invitations, len(invitationsSent))
		for _, saved := range invitationsSent {
			suite.Contains(out.Invitations, saved)
		}
	})
}

// Test_Invitations_UserInvitations тестирует получение списка приглашений направленных пользователю
func (suite *testSuite) Test_Invitations_UserInvitations() {
	suite.Run("пользователя не приглашали и потому вернется пустой список", func() {
		// Получить список приглашений
		input := ReceivedInvitationsIn{
			SubjectID: uuid.New(),
		}
		invitations, err := suite.ss.chats.ReceivedInvitations(input)
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
		input := ReceivedInvitationsIn{
			SubjectID: userID,
		}
		out, err := suite.ss.chats.ReceivedInvitations(input)
		suite.NoError(err)
		// В списке будут только приглашения, направленные пользователю
		suite.Len(out.ChatsInvitations, len(invitationsOfUser))
		for _, invitation := range out.ChatsInvitations {
			suite.Contains(invitationsOfUser, invitation)
		}
	})
}

// Test_Invitations_SendChatInvitation тестирует отправку приглашения
func (suite *testSuite) Test_Invitations_SendChatInvitation() {
	suite.Run("чат должен существовать", func() {
		// Отправить приглашение
		input := SendInvitationIn{
			SubjectID: uuid.New(),
			ChatID:    uuid.New(),
			UserID:    uuid.New(),
		}
		invitation, err := suite.ss.chats.SendInvitation(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, chatt.ErrChatNotExists)
		suite.Zero(invitation)
	})

	suite.Run("субъект должен быть участником", func() {
		// Создать чат
		chat := suite.RndChat()
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отправить приглашение
		input := SendInvitationIn{
			SubjectID: uuid.New(),
			ChatID:    chat.ID,
			UserID:    uuid.New(),
		}
		invitation, err := suite.ss.chats.SendInvitation(input)
		// Вернется ошибка, потому что субъект не является участником чата
		suite.ErrorIs(err, chatt.ErrSubjectIsNotMember)
		suite.Zero(invitation)
	})

	suite.Run("приглашаемый пользователь может не существовать", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отправить приглашение
		input := SendInvitationIn{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    uuid.New(),
		}
		out, err := suite.ss.chats.SendInvitation(input)
		suite.NoError(err)
		suite.NotZero(out)
	})

	suite.Run("приглашаемый пользователь не должен состоять в этом чате", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать участника для приглашаемого пользователя
		participantInvitating := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отправить приглашение
		input := SendInvitationIn{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    participantInvitating.UserID,
		}
		invitation, err := suite.ss.chats.SendInvitation(input)
		// Вернется ошибка, потому что приглашаемый пользователь уже является участником этого чата
		suite.ErrorIs(err, chatt.ErrParticipantExists)
		suite.Zero(invitation)
	})

	suite.Run("одновременно не может существовать несколько приглашений одного пользователя в этот чат", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать приглашаемого пользователя
		targetUserID := uuid.New()
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отправить приглашение
		input := SendInvitationIn{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    targetUserID,
		}
		invitation, err := suite.ss.chats.SendInvitation(input)
		suite.NoError(err)
		suite.Require().NotZero(invitation)
		// Отправить повторно приглашение
		invitation, err = suite.ss.chats.SendInvitation(input)
		// Вернется ошибка, потому что этот пользователь уже приглашен в чат
		suite.ErrorIs(err, chatt.ErrUserIsAlreadyInvited)
		suite.Zero(invitation)
	})

	suite.Run("любой участник может приглашать много пользователей", func() {
		// Создать чат
		chat := suite.RndChat()
		// Сохранить чат
		suite.UpsertChat(chat)
		// Создать много приглашений от разных участников
		var invitationsCreated []chatt.Invitation
		for range 5 {
			chat, err := chatt.Find(suite.rr.chats, chatt.Filter{})
			suite.Require().NoError(err)
			// Создать участника
			participant := suite.AddRndParticipant(&chat)
			// Сохранить чат
			suite.UpsertChat(chat)
			for range 5 {
				// Отправить приглашение
				input := SendInvitationIn{
					ChatID:    chat.ID,
					SubjectID: participant.UserID,
					UserID:    uuid.New(),
				}
				out, err := suite.ss.chats.SendInvitation(input)
				suite.NoError(err)
				suite.Require().NotZero(out)
				invitationsCreated = append(invitationsCreated, out.Invitation)
			}
		}
		// Получить список приглашений
		chats, err := suite.rr.chats.List(chatt.Filter{})
		suite.NoError(err)
		// В списке содержатся все созданные приглашения
		suite.Require().Len(chats, 1)
		suite.Require().Len(chats[0].Invitations, len(invitationsCreated))
		for _, createdInvitation := range invitationsCreated {
			suite.Contains(chats[0].Invitations, createdInvitation)
		}
	})
}

// Test_Invitations_AcceptInvitation тестирует принятие приглашения
func (suite *testSuite) Test_Invitations_AcceptInvitation() {
	suite.Run("приглашение должно существовать", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		p := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Принять приглашение
		input := AcceptInvitationIn{
			SubjectID:    p.UserID,
			InvitationID: uuid.New(),
		}
		err := suite.ss.chats.AcceptInvitation(input)
		suite.ErrorIs(err, ErrInvitationNotExists)
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
		input := AcceptInvitationIn{
			SubjectID:    invitation.RecipientID,
			InvitationID: invitation.ID,
		}
		err := suite.ss.chats.AcceptInvitation(input)
		suite.Require().NoError(err)
		// Получить список участников
		chats, err := suite.rr.chats.List(chatt.Filter{})
		suite.NoError(err)
		// В списке будет 3 участника: адм., приглашаемый, приглашающий
		suite.Require().Len(chats, 1)
		suite.Require().Len(chats[0].Participants, 3)
		suite.Contains(chats[0].Participants, p)
	})
}

// Test_Invitations_CancelInvitation тестирует отмену приглашения
func (suite *testSuite) Test_Invitations_CancelInvitation() {
	suite.Run("приглашение должно существовать", func() {
		// Отменить приглашение
		input := CancelInvitationIn{
			SubjectID:    uuid.New(),
			InvitationID: uuid.New(),
		}
		err := suite.ss.chats.CancelInvitation(input)
		// Вернется ошибка, потому что приглашения не существует
		suite.ErrorIs(err, ErrInvitationNotExists)
	})

	suite.Run("приглашение могут отменить только пригласивший и приглашаемый пользователи, и администратор чата", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Объявить id приглашаемого пользователя
		recipientID := uuid.New()
		// Список id тех пользователей, которые могут отменять приглашение
		cancelInvitationSubjectIDs := []uuid.UUID{
			chat.ChiefID,       // главный администратор
			participant.UserID, // пригласивший
			recipientID,        // приглашаемый
		}
		// Каждый попытается отменить приглашение
		for _, subjectUserID := range cancelInvitationSubjectIDs {
			chat, err := chatt.Find(suite.rr.chats, chatt.Filter{})
			suite.Require().NoError(err)
			// Создать приглашение
			invitation := suite.NewInvitation(participant.UserID, recipientID)
			suite.AddInvitation(&chat, invitation)
			// Сохранить чат
			suite.UpsertChat(chat)
			// Отменить приглашение
			input := CancelInvitationIn{
				SubjectID:    subjectUserID,
				InvitationID: invitation.ID,
			}
			err = suite.ss.chats.CancelInvitation(input)
			suite.NoError(err)
		}
	})

	suite.Run("другие участники не могут отменять приглашать ", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Случайный участник
		participantOther := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(participant.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отменить приглашение
		input := CancelInvitationIn{
			SubjectID:    participantOther.UserID,
			InvitationID: invitation.ID,
		}
		err := suite.ss.chats.CancelInvitation(input)
		// Вернется ошибка, потому что случайный участник не может отменять приглашение
		suite.ErrorIs(err, ErrSubjectUserNotAllowed)
	})

	suite.Run("после отмены, приглашение удаляется", func() {
		// Создать чат
		chat := suite.RndChat()
		// Создать участника
		participant := suite.AddRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.NewInvitation(participant.UserID, uuid.New())
		suite.AddInvitation(&chat, invitation)
		// Сохранить чат
		suite.UpsertChat(chat)
		// Отменить приглашение
		input := CancelInvitationIn{
			SubjectID:    invitation.SubjectID,
			InvitationID: invitation.ID,
		}
		err := suite.ss.chats.CancelInvitation(input)
		suite.Require().NoError(err)
		// Получить список приглашений
		chats, err := suite.rr.chats.List(chatt.Filter{})
		suite.NoError(err)
		suite.Require().Len(chats, 1)
		suite.Empty(chats[0].Invitations)
	})
}
