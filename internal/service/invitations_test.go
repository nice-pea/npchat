package service

import (
	"testing"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

// Test_ChatInvitationsInput_Validate тестирует валидацию входящих параметров
func Test_ChatInvitationsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := ChatInvitationsIn{
			SubjectID: id,
			ChatID:    id,
		}
		return input.Validate()
	})
}

// Test_Invitations_ChatInvitations тестирует получение списка приглашений
func (suite *servicesTestSuite) Test_Invitations_ChatInvitations() {
	suite.Run("чат должен существовать", func() {
		// Получить список приглашений
		input := ChatInvitationsIn{
			SubjectID: uuid.NewString(),
			ChatID:    uuid.NewString(),
		}
		invitations, err := suite.ss.chats.ChatInvitations(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
		suite.Empty(invitations)
	})

	suite.Run("субъект должен быть участником чата", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Получить список приглашений
		input := ChatInvitationsIn{
			ChatID:    chat.ID,
			SubjectID: uuid.NewString(),
		}
		invitations, err := suite.ss.chats.ChatInvitations(input)
		// Вернется ошибка, потому что пользователь не участник чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
		suite.Empty(invitations)
	})

	suite.Run("пустой список из чата без приглашений", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Получить список приглашений
		input := ChatInvitationsIn{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
		}
		invitations, err := suite.ss.chats.ChatInvitations(input)
		suite.NoError(err)
		suite.Empty(invitations)
	})

	suite.Run("субъект не администратор чата и видит только отправленные им приглашения", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		participant := suite.addRndParticipant(&chat)
		// Создать приглашения отправленные участником
		subjectInvitations := make([]chatt.Invitation, 3)
		for i := range subjectInvitations {
			subjectInvitations[i] = suite.newInvitation(participant.UserID, uuid.NewString())
			suite.addInvitation(&chat, subjectInvitations[i])
		}
		// Создать приглашения отправленные какими-то другими пользователями
		for range 3 {
			p := suite.addRndParticipant(&chat)
			i := suite.newInvitation(p.UserID, uuid.NewString())
			suite.addInvitation(&chat, i)
		}
		// Получить список приглашений
		input := ChatInvitationsIn{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
		}
		out, err := suite.ss.chats.ChatInvitations(input)
		suite.Require().NoError(err)
		// В списке будут приглашения отправленные участником
		suite.Require().Len(out.Invitations, len(subjectInvitations))
		for i, subjectInvitation := range subjectInvitations {
			suite.Equal(subjectInvitation, out.Invitations[i])
		}
	})

	suite.Run("субъект является администратором чата и видит все отправленные приглашения в чат", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать приглашения отправленные какими-то пользователями
		invitationsSent := make([]chatt.Invitation, 3)
		for i := range invitationsSent {
			p := suite.addRndParticipant(&chat)
			invitationsSent[i] = suite.newInvitation(p.UserID, uuid.NewString())
			suite.addInvitation(&chat, invitationsSent[i])
		}
		// Получить список приглашений
		input := ChatInvitationsIn{
			SubjectID: chat.ChiefID,
			ChatID:    chat.ID,
		}
		out, err := suite.ss.chats.ChatInvitations(input)
		suite.Require().NoError(err)
		// В списке будут приглашения все приглашения
		suite.Require().Len(out, len(invitationsSent))
		for _, saved := range invitationsSent {
			suite.Contains(out.Invitations, saved)
		}
	})
}

// Test_UserInvitationsInput_Validate тестирует валидацию входящих параметров
func Test_UserInvitationsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := ReceivedInvitationsIn{
			SubjectID: id,
		}
		return input.Validate()
	})
}

// Test_Invitations_UserInvitations тестирует получение списка приглашений направленных пользователю
func (suite *servicesTestSuite) Test_Invitations_UserInvitations() {
	suite.Run("пользователя не приглашали и потому вернется пустой список", func() {
		// Получить список приглашений
		input := ReceivedInvitationsIn{
			SubjectID: uuid.NewString(),
		}
		invitations, err := suite.ss.chats.ReceivedInvitations(input)
		suite.NoError(err)
		suite.Empty(invitations)
	})

	suite.Run("из многих приглашений несколько направлено пользователю и потому вернутся только они", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// ID пользователя
		userID := uuid.NewString()
		// Создать несколько приглашений направленных пользователю
		invitationsOfUser := make([]chatt.Invitation, 5)
		for i := range invitationsOfUser {
			p := suite.addRndParticipant(&chat)
			invitationsOfUser[i] = suite.newInvitation(p.UserID, userID)
		}
		// Создать несколько приглашений направленных каким-то другим пользователям
		for range 10 {
			p := suite.addRndParticipant(&chat)
			i := suite.newInvitation(p.UserID, uuid.NewString())
			suite.addInvitation(&chat, i)
		}
		// Получить список приглашений
		input := ReceivedInvitationsIn{
			SubjectID: userID,
		}
		out, err := suite.ss.chats.ReceivedInvitations(input)
		suite.NoError(err)
		// В списке будут только приглашения направленные пользователю
		suite.Require().Len(out, len(invitationsOfUser))
		for _, invitation := range invitationsOfUser {
			suite.Contains(out.ChatsInvitations, invitation)
		}
	})
}

// Test_SendChatInvitationInput_Validate тестирует валидацию входящих параметров
func Test_SendChatInvitationInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := SendInvitationIn{
			SubjectID: id,
			ChatID:    id,
			UserID:    id,
		}
		return input.Validate()
	})
}

// Test_Invitations_SendChatInvitation тестирует отправку приглашения
func (suite *servicesTestSuite) Test_Invitations_SendChatInvitation() {
	suite.Run("чат должен существовать", func() {
		// Отправить приглашение
		input := SendInvitationIn{
			SubjectID: uuid.NewString(),
			ChatID:    uuid.NewString(),
			UserID:    uuid.NewString(),
		}
		invitation, err := suite.ss.chats.SendInvitation(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
		suite.Zero(invitation)
	})

	suite.Run("субъект должен быть участником", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Отправить приглашение
		input := SendInvitationIn{
			SubjectID: uuid.NewString(),
			ChatID:    chat.ID,
			UserID:    uuid.NewString(),
		}
		invitation, err := suite.ss.chats.SendInvitation(input)
		// Вернется ошибка, потому что субъект не является участником чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
		suite.Zero(invitation)
	})

	suite.Run("приглашаемый пользователь должен существовать", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		participant := suite.addRndParticipant(&chat)
		// Отправить приглашение
		input := SendInvitationIn{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    uuid.NewString(),
		}
		invitation, err := suite.ss.chats.SendInvitation(input)
		// Вернется ошибка, потому что приглашаемого пользователя не существует
		suite.ErrorIs(err, ErrUserNotExists)
		suite.Zero(invitation)
	})

	suite.Run("приглашаемый пользователь не должен состоять в этом чате", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		participant := suite.addRndParticipant(&chat)
		// Создать участника для приглашаемого пользователя
		participantInvitating := suite.addRndParticipant(&chat)
		// Отправить приглашение
		input := SendInvitationIn{
			ChatID:    chat.ID,
			SubjectID: participant.UserID,
			UserID:    participantInvitating.UserID,
		}
		invitation, err := suite.ss.chats.SendInvitation(input)
		// Вернется ошибка, потому что приглашаемый пользователь уже является участником этого чата
		suite.ErrorIs(err, ErrUserIsAlreadyInChat)
		suite.Zero(invitation)
	})

	suite.Run("одновременно не может существовать несколько приглашений одного пользователя в этот чат", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		participant := suite.addRndParticipant(&chat)
		// Создать приглашаемого пользователя
		targetUserID := uuid.NewString()
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
		suite.ErrorIs(err, ErrUserIsAlreadyInvited)
		suite.Zero(invitation)
	})

	suite.Run("любой участник может приглашать много пользователей", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать много приглашений от разных участников
		var invitationsCreated []chatt.Invitation
		for range 5 {
			// Создать участника
			participant := suite.addRndParticipant(&chat)
			for range 5 {
				// Отправить приглашение
				input := SendInvitationIn{
					ChatID:    chat.ID,
					SubjectID: participant.UserID,
					UserID:    uuid.NewString(),
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

// Test_AcceptInvitationInput_Validate тестирует валидацию входящих параметров
func Test_AcceptInvitationInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		inp := AcceptInvitationIn{
			SubjectID:    id,
			InvitationID: id,
		}
		return inp.Validate()
	})
}

// Test_Invitations_AcceptInvitation тестирует принятие приглашения
func (suite *servicesTestSuite) Test_Invitations_AcceptInvitation() {
	suite.Run("приглашение должно существовать", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		p := suite.addRndParticipant(&chat)
		// Принять приглашение
		input := AcceptInvitationIn{
			SubjectID:    p.UserID,
			InvitationID: uuid.NewString(),
		}
		err := suite.ss.chats.AcceptInvitation(input)
		suite.ErrorIs(err, ErrInvitationNotExists)
	})

	suite.Run("приняв приглашение, пользователь становится участником чата", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		p := suite.addRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.newInvitation(p.UserID, uuid.NewString())
		suite.addInvitation(&chat, invitation)
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
		// В списке будет только один участник, который принял приглашение
		suite.Require().Len(chats, 1)
		suite.Require().Len(chats[0].Participants, 1)
		suite.Equal(p, chats[0].Participants[0])
	})
}

// Test_CancelInvitationInput_Validate тестирует валидацию входящих параметров
func Test_CancelInvitationInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := CancelInvitationIn{
			SubjectID:    id,
			InvitationID: id,
		}
		return input.Validate()
	})
}

// Test_Invitations_CancelInvitation тестирует отмену приглашения
func (suite *servicesTestSuite) Test_Invitations_CancelInvitation() {
	suite.Run("приглашение должно существовать", func() {
		// Отменить приглашение
		input := CancelInvitationIn{
			SubjectID:    uuid.NewString(),
			InvitationID: uuid.NewString(),
		}
		err := suite.ss.chats.CancelInvitation(input)
		// Вернется ошибка, потому что приглашения не существует
		suite.ErrorIs(err, ErrInvitationNotExists)
	})

	suite.Run("приглашение могут отменить только пригласивший и приглашаемый пользователи, и администратор чата", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		participant := suite.addRndParticipant(&chat)
		// Объявить id приглашаемого пользователя
		recipientID := uuid.NewString()
		// Список id тех пользователей, которые могут отменять приглашение
		cancelInvitationSubjectIDs := []string{
			chat.ChiefID,       // главный администратор
			participant.UserID, // пригласивший
			recipientID,        // приглашаемый
		}
		// Каждый попытается отменить приглашение
		for _, subjectUserID := range cancelInvitationSubjectIDs {
			// Создать приглашение
			invitation := suite.newInvitation(participant.UserID, recipientID)
			suite.addInvitation(&chat, invitation)
			// Отменить приглашение
			input := CancelInvitationIn{
				SubjectID:    subjectUserID,
				InvitationID: invitation.ID,
			}
			err := suite.ss.chats.CancelInvitation(input)
			suite.NoError(err)
		}
	})

	suite.Run("другие участники не могут отменять приглашать ", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		participant := suite.addRndParticipant(&chat)
		// Случайный участник
		participantOther := suite.addRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.newInvitation(participant.UserID, uuid.NewString())
		suite.addInvitation(&chat, invitation)
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
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		participant := suite.addRndParticipant(&chat)
		// Создать приглашение
		invitation := suite.newInvitation(participant.UserID, uuid.NewString())
		suite.addInvitation(&chat, invitation)
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
