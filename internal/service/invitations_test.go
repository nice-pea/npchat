package service

import (
	"testing"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
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
		// Создать приглашения отправленные участником
		subjectInvitations := make([]domain.Invitation, 3)
		for i := range subjectInvitations {
			subjectInvitations[i] = suite.saveInvitation(domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member.UserID,
				UserID:        uuid.NewString(),
				ChatID:        chat.ID,
			})
		}
		// Создать приглашения отправленные какими-то другими пользователями
		for range 3 {
			suite.saveInvitation(domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: uuid.NewString(),
				UserID:        uuid.NewString(),
				ChatID:        chat.ID,
			})
		}
		// Получить список приглашений
		input := ChatInvitationsIn{
			ChatID:    chat.ID,
			SubjectID: member.UserID,
		}
		invitationsFromService, err := suite.ss.invitations.ChatInvitations(input)
		suite.Require().NoError(err)
		// В списке будут приглашения отправленные участником
		if suite.Len(invitationsFromService, len(subjectInvitations)) {
			for i, subjectInvitation := range subjectInvitations {
				suite.Equal(subjectInvitation, invitationsFromService[i])
			}
		}
	})

	suite.Run("субъект является администратором чата и видит все отправленные приглашения в чат", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать приглашения отправленные какими-то пользователями
		invitationsSaved := make([]domain.Invitation, 4)
		for i := range invitationsSaved {
			invitationsSaved[i] = suite.saveInvitation(domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: chat.ID,
				UserID: uuid.NewString(),
			})
		}
		// Получить список приглашений
		input := ChatInvitationsIn{
			SubjectID: member.UserID,
			ChatID:    chat.ID,
		}
		invitationFromService, err := suite.ss.invitations.ChatInvitations(input)
		suite.Require().NoError(err)
		// В списке будут приглашения все приглашения
		if suite.Len(invitationFromService, len(invitationsSaved)) {
			for i, saved := range invitationsSaved {
				suite.Equal(saved, invitationFromService[i])
			}
		}
	})
}

// Test_UserInvitationsInput_Validate тестирует валидацию входящих параметров
func Test_UserInvitationsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := ReceivedInvitationsIn{
			SubjectID: id,
			UserID:    id,
		}
		return input.Validate()
	})
}

// Test_Invitations_UserInvitations тестирует получение списка приглашений направленных пользователю
func (suite *servicesTestSuite) Test_Invitations_UserInvitations() {
	suite.Run("пользователь должен существовать", func() {
		id := uuid.NewString()
		// Получить список приглашений
		input := ReceivedInvitationsIn{
			SubjectID: id,
			UserID:    id,
		}
		invitations, err := suite.ss.invitations.ReceivedInvitations(input)
		// Вернется ошибка, потому что пользователя не существует
		suite.ErrorIs(err, ErrUserNotExists)
		suite.Empty(invitations)
	})

	suite.Run("пользователь может просматривать только свои приглашения", func() {
		// Получить список приглашений
		input := ReceivedInvitationsIn{
			SubjectID: uuid.NewString(),
			UserID:    uuid.NewString(),
		}
		invitations, err := suite.ss.invitations.ReceivedInvitations(input)
		// Вернется ошибка, потому что пользователь пытается просмотреть чужие приглашения
		suite.ErrorIs(err, ErrUnauthorizedInvitationsView)
		suite.Empty(invitations)
	})

	suite.Run("пользователя не приглашали и потому вернется пустой список", func() {
		// Создать пользователя
		user := suite.saveUser(domain.User{
			ID: uuid.NewString(),
		})
		// Получить список приглашений
		input := ReceivedInvitationsIn{
			SubjectID: user.ID,
			UserID:    user.ID,
		}
		invitations, err := suite.ss.invitations.ReceivedInvitations(input)
		suite.NoError(err)
		suite.Empty(invitations)
	})

	suite.Run("из многих приглашений несколько направлено пользователю и потому вернутся только они", func() {
		// Создать пользователя
		user := suite.saveUser(domain.User{
			ID: uuid.NewString(),
		})
		// Создать несколько приглашений направленных пользователю
		invitationsOfUser := make([]domain.Invitation, 5)
		for i := range invitationsOfUser {
			invitationsOfUser[i] = suite.saveInvitation(domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
				UserID: user.ID,
			})
		}
		// Создать несколько приглашений направленных каким-то другим пользователям
		for range 10 {
			suite.saveInvitation(domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
				UserID: uuid.NewString(),
			})
		}
		// Получить список приглашений
		input := ReceivedInvitationsIn{
			SubjectID: user.ID,
			UserID:    user.ID,
		}
		invitationsFromRepo, err := suite.ss.invitations.ReceivedInvitations(input)
		suite.NoError(err)
		// В списке будут только приглашения направленные пользователю
		if suite.Len(invitationsFromRepo, len(invitationsOfUser)) {
			for i, invitation := range invitationsOfUser {
				suite.Equal(invitation.ID, invitationsFromRepo[i].ID)
				suite.Equal(invitation.ChatID, invitationsFromRepo[i].ChatID)
				suite.Equal(invitation.UserID, invitationsFromRepo[i].UserID)
			}
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
		invitation, err := suite.ss.invitations.SendInvitation(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
		suite.Zero(invitation)
	})

	suite.Run("субъект должен быть участником", func() {
		// Создать чат
		chat := suite.upsertChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Отправить приглашение
		input := SendInvitationIn{
			SubjectID: uuid.NewString(),
			ChatID:    chat.ID,
			UserID:    uuid.NewString(),
		}
		invitation, err := suite.ss.invitations.SendInvitation(input)
		// Вернется ошибка, потому что субъект не является участником чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
		suite.Zero(invitation)
	})

	suite.Run("приглашаемый пользователь должен существовать", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Отправить приглашение
		input := SendInvitationIn{
			ChatID:    chat.ID,
			SubjectID: member.UserID,
			UserID:    uuid.NewString(),
		}
		invitation, err := suite.ss.invitations.SendInvitation(input)
		// Вернется ошибка, потому что приглашаемого пользователя не существует
		suite.ErrorIs(err, ErrUserNotExists)
		suite.Zero(invitation)
	})

	suite.Run("приглашаемый пользователь не должен состоять в этом чате", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		subjectMember := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Создать приглашаемого пользователя
		targetUser := suite.saveUser(domain.User{
			ID: uuid.NewString(),
		})
		// Создать участника для приглашаемого пользователя
		suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: targetUser.ID,
			ChatID: chat.ID,
		})
		// Отправить приглашение
		input := SendInvitationIn{
			ChatID:    chat.ID,
			SubjectID: subjectMember.UserID,
			UserID:    targetUser.ID,
		}
		invitation, err := suite.ss.invitations.SendInvitation(input)
		// Вернется ошибка, потому что приглашаемый пользователь уже является участником этого чата
		suite.ErrorIs(err, ErrUserIsAlreadyInChat)
		suite.Zero(invitation)
	})

	suite.Run("одновременно не может существовать несколько приглашений одного пользователя в этот чат", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать участника
		subjectMember := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Создать приглашаемого пользователя
		targetUser := suite.saveUser(domain.User{
			ID: uuid.NewString(),
		})
		// Отправить приглашение
		input := SendInvitationIn{
			ChatID:    chat.ID,
			SubjectID: subjectMember.UserID,
			UserID:    targetUser.ID,
		}
		invitation, err := suite.ss.invitations.SendInvitation(input)
		suite.NoError(err)
		suite.Require().NotZero(invitation)
		// Отправить повторно приглашение
		invitation, err = suite.ss.invitations.SendInvitation(input)
		// Вернется ошибка, потому что этот пользователь уже приглашен в чат
		suite.ErrorIs(err, ErrUserIsAlreadyInvited)
		suite.Zero(invitation)
	})

	suite.Run("любой участник может приглашать много пользователей", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать много приглашений от разных участников
		var createdInvitations []domain.Invitation
		for range 5 {
			// Создать участника
			subjectMember := suite.saveMember(domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			})
			for range 5 {
				// Создать приглашаемого пользователя
				targetUser := suite.saveUser(domain.User{
					ID: uuid.NewString(),
				})
				// Отправить приглашение
				input := SendInvitationIn{
					ChatID:    chat.ID,
					SubjectID: subjectMember.UserID,
					UserID:    targetUser.ID,
				}
				invitation, err := suite.ss.invitations.SendInvitation(input)
				suite.NoError(err)
				suite.Require().NotZero(invitation)
				createdInvitations = append(createdInvitations, invitation)
			}
		}
		// Получить список приглашений
		invitationsFromRepo, err := suite.ss.invitations.InvitationsRepo.List(domain.InvitationsFilter{})
		suite.NoError(err)
		// В списке содержатся все созданные приглашения
		suite.Require().Len(invitationsFromRepo, len(createdInvitations))
		for i, createdInvitation := range createdInvitations {
			suite.Equal(createdInvitation.ChatID, invitationsFromRepo[i].ChatID)
			suite.Equal(createdInvitation.SubjectUserID, invitationsFromRepo[i].SubjectUserID)
			suite.Equal(createdInvitation.UserID, invitationsFromRepo[i].UserID)
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
		// Принять приглашение
		input := AcceptInvitationIn{
			SubjectID:    uuid.NewString(),
			InvitationID: uuid.NewString(),
		}
		err := suite.ss.invitations.AcceptInvitation(input)
		suite.ErrorIs(err, ErrInvitationNotExists)
	})

	suite.Run("пользователь должен существовать", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		// Создать приглашение
		invitation := suite.saveInvitation(domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
			ChatID:        chat.ID,
		})
		// Принять приглашение
		input := AcceptInvitationIn{
			SubjectID:    invitation.UserID,
			InvitationID: invitation.ID,
		}
		err := suite.ss.invitations.AcceptInvitation(input)
		suite.ErrorIs(err, ErrUserNotExists)
	})

	suite.Run("приглашение быть направлено пользователю", func() {
		// Создать чат
		user := suite.saveUser(domain.User{
			ID: uuid.NewString(),
		})
		// Создать приглашение
		invitation := suite.saveInvitation(domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
		})
		// Принять приглашение
		input := AcceptInvitationIn{
			SubjectID:    user.ID,
			InvitationID: invitation.ID,
		}
		err := suite.ss.invitations.AcceptInvitation(input)
		suite.ErrorIs(err, ErrSubjectUserNotAllowed)
	})

	suite.Run("приняв приглашение, пользователь становится участником чата", func() {
		// Создать чат
		chat := suite.upsertChat(suite.rndChat())
		user := suite.saveUser(domain.User{
			ID: uuid.NewString(),
		})
		// Создать приглашение
		invitation := suite.saveInvitation(domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        user.ID,
			ChatID:        chat.ID,
		})
		// Принять приглашение
		input := AcceptInvitationIn{
			SubjectID:    user.ID,
			InvitationID: invitation.ID,
		}
		err := suite.ss.invitations.AcceptInvitation(input)
		suite.Require().NoError(err)
		// Получить список участников
		members, err := suite.ss.invitations.MembersRepo.List(domain.MembersFilter{})
		suite.NoError(err)
		// В списке будет только один участник, который принял приглашение
		suite.Require().Len(members, 1)
		suite.Equal(user.ID, members[0].UserID)
		suite.Equal(chat.ID, members[0].ChatID)
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
		err := suite.ss.invitations.CancelInvitation(input)
		// Вернется ошибка, потому что приглашения не существует
		suite.ErrorIs(err, ErrInvitationNotExists)
	})

	suite.Run("приглашение могут отменить только пригласивший и приглашаемый пользователи, и администратор чата", func() {
		// Создать чат
		chat := suite.upsertChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Создать участника
		subjectMember := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Объявить id приглашаемого пользователя
		userID := uuid.NewString()

		// Список id тех пользователей, которые могут отменять приглашение
		cancelInvitationSubjectIDs := []string{
			chat.ChiefUserID,     // главный администратор
			subjectMember.UserID, // пригласивший
			userID,               // приглашаемый
		}
		// Каждый попытается отменить приглашение
		for _, subjectUserID := range cancelInvitationSubjectIDs {
			// Создать приглашение
			invitation := suite.saveInvitation(domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: subjectMember.UserID,
				UserID:        userID,
				ChatID:        chat.ID,
			})
			// Отменить приглашение
			input := CancelInvitationIn{
				SubjectID:    subjectUserID,
				InvitationID: invitation.ID,
			}
			err := suite.ss.invitations.CancelInvitation(input)
			suite.NoError(err)
		}
	})

	suite.Run("другие участники не могут отменять приглашать ", func() {
		// Случайный участник
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: uuid.NewString(),
		})
		// Создать приглашение
		invitation := suite.saveInvitation(domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
			ChatID:        member.ChatID,
		})
		// Отменить приглашение
		input := CancelInvitationIn{
			SubjectID:    member.UserID,
			InvitationID: invitation.ID,
		}
		err := suite.ss.invitations.CancelInvitation(input)
		// Вернется ошибка, потому что случайный участник не может отменять приглашение
		suite.ErrorIs(err, ErrSubjectUserNotAllowed)
	})

	suite.Run("после отмены, приглашение удаляется", func() {
		// Создать участника
		subjectMember := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
		})
		// Создать приглашение
		invitation := suite.saveInvitation(domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: subjectMember.UserID,
			UserID:        uuid.NewString(),
		})
		// Отменить приглашение
		input := CancelInvitationIn{
			SubjectID:    invitation.SubjectUserID,
			InvitationID: invitation.ID,
		}
		err := suite.ss.invitations.CancelInvitation(input)
		suite.Require().NoError(err)
		// Получить список приглашений
		invitations, err := suite.ss.invitations.InvitationsRepo.List(domain.InvitationsFilter{})
		suite.NoError(err)
		suite.Empty(invitations)
	})
}
