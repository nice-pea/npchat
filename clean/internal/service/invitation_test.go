package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

func (suite *servicesTestSuite) saveInvitation(invitation domain.Invitation) domain.Invitation {
	err := suite.invitationsService.InvitationsRepo.Save(invitation)
	suite.Require().NoError(err)

	return invitation
}

func (suite *servicesTestSuite) saveUser(user domain.User) domain.User {
	err := suite.invitationsService.UsersRepo.Save(user)
	suite.Require().NoError(err)

	return user
}

func Test_ChatInvitationsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := ChatInvitationsInput{
			SubjectUserID: id,
			ChatID:        id,
		}
		return input.Validate()
	})
}

func (suite *servicesTestSuite) Test_Invitations_ChatInvitations() {
	suite.Run("чат должен существовать", func() {
		input := ChatInvitationsInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
		}
		invitations, err := suite.invitationsService.ChatInvitations(input)
		suite.ErrorIs(err, ErrChatNotExists)
		suite.Empty(invitations)
	})
	suite.Run("субъект должен быть участником чата", func() {
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		input := ChatInvitationsInput{
			ChatID:        chat.ID,
			SubjectUserID: uuid.NewString(),
		}
		invitations, err := suite.invitationsService.ChatInvitations(input)
		suite.Error(err)
		suite.Empty(invitations)
	})
	suite.Run("пустой список из чата без приглашений", func() {
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		})
		input := ChatInvitationsInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		invitations, err := suite.invitationsService.ChatInvitations(input)
		suite.NoError(err)
		suite.Empty(invitations)
	})
	suite.Run("субъект не администратор чата и видит только отправленные им приглашения", func() {
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})

		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})

		subjectInvitations := make([]domain.Invitation, 3)
		for i := range subjectInvitations {
			subjectInvitations[i] = suite.saveInvitation(domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member.UserID,
				UserID:        uuid.NewString(),
				ChatID:        chat.ID,
			})
		}

		for range 3 {
			suite.saveInvitation(domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: uuid.NewString(),
				UserID:        uuid.NewString(),
				ChatID:        chat.ID,
			})
		}

		input := ChatInvitationsInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
		}
		invitationsFromService, err := suite.invitationsService.ChatInvitations(input)
		suite.Require().NoError(err)

		if suite.Len(invitationsFromService, len(subjectInvitations)) {
			for i, subjectInvitation := range subjectInvitations {
				suite.Equal(subjectInvitation, invitationsFromService[i])
			}
		}
	})
	suite.Run("субъект является администратором чата и видит все отправленные приглашения в чат", func() {
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})

		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		})

		invitationsSaved := make([]domain.Invitation, 4)
		for i := range invitationsSaved {
			invitationsSaved[i] = suite.saveInvitation(domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: chat.ID,
			})
		}

		input := ChatInvitationsInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		invitationFromService, err := suite.invitationsService.ChatInvitations(input)
		suite.Require().NoError(err)

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
		input := UserInvitationsInput{
			SubjectUserID: id,
			UserID:        id,
		}
		return input.Validate()
	})
}

// Test_Invitations_UserInvitations тестирование функции UserInvitations
func (suite *servicesTestSuite) Test_Invitations_UserInvitations() {
	suite.Run("пользователь должен существовать", func() {
		id := uuid.NewString()
		input := UserInvitationsInput{
			SubjectUserID: id,
			UserID:        id,
		}
		invitations, err := suite.invitationsService.UserInvitations(input)
		suite.ErrorIs(err, ErrUserNotExists)
		suite.Empty(invitations)
	})
	suite.Run("пользователь может просматривать только свои приглашения", func() {
		input := UserInvitationsInput{
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		invitations, err := suite.invitationsService.UserInvitations(input)
		suite.ErrorIs(err, ErrUnauthorizedInvitationsView)
		suite.Empty(invitations)
	})
	suite.Run("пользователя не приглашали и потому вернется пустой список", func() {
		user := suite.saveUser(domain.User{
			ID: uuid.NewString(),
		})
		input := UserInvitationsInput{
			SubjectUserID: user.ID,
			UserID:        user.ID,
		}
		invitations, err := suite.invitationsService.UserInvitations(input)
		suite.NoError(err)
		suite.Empty(invitations)
	})
	suite.Run("из многих приглашений несколько направлено пользователю и потому вернутся только они", func() {
		user := suite.saveUser(domain.User{
			ID: uuid.NewString(),
		})

		invitationsOfUser := make([]domain.Invitation, 5)
		for i := range invitationsOfUser {
			invitationsOfUser[i] = suite.saveInvitation(domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
				UserID: user.ID,
			})
		}

		for range 10 {
			suite.saveInvitation(domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
				UserID: uuid.NewString(),
			})
		}

		input := UserInvitationsInput{
			SubjectUserID: user.ID,
			UserID:        user.ID,
		}
		invitationsFromRepo, err := suite.invitationsService.UserInvitations(input)
		suite.NoError(err)

		if suite.Len(invitationsFromRepo, len(invitationsOfUser)) {
			for i, invitation := range invitationsOfUser {
				suite.Equal(invitation.ID, invitationsFromRepo[i].ID)
				suite.Equal(invitation.ChatID, invitationsFromRepo[i].ChatID)
				suite.Equal(invitation.UserID, invitationsFromRepo[i].UserID)
			}
		}
	})
}

func Test_SendChatInvitationInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := SendInvitationInput{
			SubjectUserID: id,
			ChatID:        id,
			UserID:        id,
		}
		return input.Validate()
	})
}

func Test_Invitations_SendChatInvitation(t *testing.T) {
	suite.Run("участник отправляющий приглашения должен состоять в чате", func() {
		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		subjectUser := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(subjectUser)
		suite.NoError(err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: subjectUser.ID,
		}
		err = suite.invitationsService.MembersRepo.Save(member)
		suite.NoError(err)

		targetUser := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(targetUser)
		suite.NoError(err)

		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser.ID,
		}
		err = suite.invitationsService.SendInvitation(input)
		suite.ErrorIs(err, ErrSubjectUserIsNotMember)
	})
	suite.Run("UserID не должен состоять в чате ChatID", func() {
		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = suite.invitationsService.MembersRepo.Save(member)
		suite.NoError(err)

		targetUser := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(targetUser)
		suite.NoError(err)

		targetMember := domain.Member{
			ID:     uuid.NewString(),
			UserID: targetUser.ID,
			ChatID: chat.ID,
		}
		err = suite.invitationsService.MembersRepo.Save(targetMember)
		suite.NoError(err)

		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser.ID,
		}
		err = suite.invitationsService.SendInvitation(input)
		suite.ErrorIs(err, ErrUserAlreadyInChat)
	})
	suite.Run("приглашать участников могут все члены чата", func() {
		suite.Run("админситратор", func() {
			suite.invitationsService := newInvitationsService(t)
			chatId := uuid.NewString()
			userId := uuid.NewString()

			chief := domain.Member{
				ID:     uuid.NewString(),
				UserID: userId,
				ChatID: chatId,
			}
			err := suite.invitationsService.MembersRepo.Save(chief)
			suite.NoError(err)
			chat := domain.Chat{
				ID:          chatId,
				ChiefUserID: userId,
			}
			err = suite.invitationsService.ChatsRepo.Save(chat)
			suite.NoError(err)

			targetUser := domain.User{
				ID: uuid.NewString(),
			}
			err = suite.invitationsService.UsersRepo.Save(targetUser)
			suite.NoError(err)

			input := SendInvitationInput{
				ChatID:        chat.ID,
				SubjectUserID: chief.UserID,
				UserID:        targetUser.ID,
			}
			err = suite.invitationsService.SendInvitation(input)
			suite.NoError(err)
		})
		suite.Run("обычный участник чата", func() {
			suite.invitationsService := newInvitationsService(t)

			chat := domain.Chat{
				ID: uuid.NewString(),
			}
			err := suite.invitationsService.ChatsRepo.Save(chat)
			suite.NoError(err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = suite.invitationsService.MembersRepo.Save(member)
			suite.NoError(err)

			targetUser := domain.User{
				ID: uuid.NewString(),
			}
			err = suite.invitationsService.UsersRepo.Save(targetUser)
			suite.NoError(err)

			input := SendInvitationInput{
				ChatID:        chat.ID,
				SubjectUserID: member.UserID,
				UserID:        targetUser.ID,
			}
			err = suite.invitationsService.SendInvitation(input)
			suite.NoError(err)
		})
	})
	suite.Run("UserID должен существовать", func() {
		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = suite.invitationsService.MembersRepo.Save(member)
		suite.NoError(err)

		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        uuid.NewString(),
		}
		err = suite.invitationsService.SendInvitation(input)
		suite.ErrorIs(err, ErrUserNotExists)
	})
	suite.Run("ChatID должен существовать", func() {

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		err := suite.invitationsService.MembersRepo.Save(member)
		suite.NoError(err)

		targetUser := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(targetUser)
		suite.NoError(err)

		input := SendInvitationInput{
			ChatID:        uuid.NewString(),
			SubjectUserID: member.ID,
			UserID:        targetUser.ID,
		}
		err = suite.invitationsService.SendInvitation(input)
		suite.ErrorIs(err, ErrChatNotExists)
	})
	suite.Run("UserID нельзя приглашать более 1 раза", func() {

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = suite.invitationsService.MembersRepo.Save(member)
		suite.NoError(err)

		targetUser := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(targetUser)
		suite.NoError(err)

		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser.ID,
		}
		err = suite.invitationsService.SendInvitation(input)
		suite.NoError(err)

		err = suite.invitationsService.SendInvitation(input)
		suite.ErrorIs(err, ErrUserAlreadyInviteInChat)
	})
	suite.Run("можно приглашать больее 1 раза разных пользователей", func() {

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = suite.invitationsService.MembersRepo.Save(member)
		suite.NoError(err)

		targetUser1 := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(targetUser1)
		suite.NoError(err)

		input1 := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser1.ID,
		}
		err = suite.invitationsService.SendInvitation(input1)
		suite.NoError(err)

		targetUser2 := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(targetUser2)
		suite.NoError(err)

		input2 := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser2.ID,
		}

		err = suite.invitationsService.SendInvitation(input2)
		suite.NoError(err)

		invsRepo, err := suite.invitationsService.InvitationsRepo.List(domain.InvitationsFilter{})
		suite.NoError(err)
		suite.Len(invsRepo, 2)
		for i, invInput := range []SendInvitationInput{input1, input2} {
			suite.Equal(invInput.ChatID, invsRepo[i].ChatID)
			suite.Equal(invInput.SubjectUserID, invsRepo[i].SubjectUserID)
			suite.Equal(invInput.UserID, invsRepo[i].UserID)
		}
	})
}

func Test_AcceptInvitationInput_Validate() {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		inp := AcceptInvitationInput{
			SubjectUserID: id,
			ChatID:        id,
		}
		return inp.Validate()
	})
}

func Test_Invitations_AcceptInvitation() {
	suite.Run("принятие не существующего приглашения", func() {

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(user)
		suite.NoError(err)

		input := AcceptInvitationInput{
			SubjectUserID: user.ID,
			ChatID:        chat.ID,
		}
		err = suite.invitationsService.AcceptInvitation(input)
		suite.ErrorIs(err, ErrInvitationNotExists)

		members, err := suite.invitationsService.MembersRepo.List(domain.MembersFilter{})
		suite.NoError(err)
		suite.Len(members, 0)
	})
	suite.Run("после принятия существующего приглашения, пользователь становится участником чата", func() {

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(user)
		suite.NoError(err)

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        user.ID,
			ChatID:        chat.ID,
		}
		err = suite.invitationsService.InvitationsRepo.Save(invitation)
		suite.NoError(err)

		input := AcceptInvitationInput{
			SubjectUserID: user.ID,
			ChatID:        chat.ID,
		}
		err = suite.invitationsService.AcceptInvitation(input)
		suite.NoError(err)

		members, err := suite.invitationsService.MembersRepo.List(domain.MembersFilter{})
		suite.NoError(err)
		suite.Len(members, 1)
		suite.Equal(user.ID, members[0].UserID)
		suite.Equal(chat.ID, members[0].ChatID)
	})
	suite.Run("принятие существующего приглашения в несуществющий чат", func() {

		user := domain.User{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.UsersRepo.Save(user)
		suite.NoError(err)

		chatId := uuid.NewString()

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        user.ID,
			ChatID:        chatId,
		}
		err = suite.invitationsService.InvitationsRepo.Save(invitation)
		suite.NoError(err)

		input := AcceptInvitationInput{
			SubjectUserID: user.ID,
			ChatID:        chatId,
		}
		err = suite.invitationsService.AcceptInvitation(input)
		suite.ErrorIs(err, ErrChatNotExists)
	})
	suite.Run("пользователя не существует", func() {

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
			ChatID:        chat.ID,
		}
		err = suite.invitationsService.InvitationsRepo.Save(invitation)
		suite.NoError(err)

		input := AcceptInvitationInput{
			SubjectUserID: invitation.UserID,
			ChatID:        chat.ID,
		}
		err = suite.invitationsService.AcceptInvitation(input)
		suite.ErrorIs(err, ErrUserNotExists)

		members, err := suite.invitationsService.MembersRepo.List(domain.MembersFilter{})
		suite.NoError(err)
		suite.Len(members, 0)
	})
}

func Test_CancelInvitationInput_Validate() {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := CancelInvitationInput{
			SubjectUserID: id,
			UserID:        id,
			ChatID:        id,
		}
		return input.Validate()
	})
}

func Test_Invitations_CancelInvitation() {
	suite.Run("отменить не существующее приглашение", func() {

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(user)
		suite.NoError(err)

		input := CancelInvitationInput{
			SubjectUserID: user.ID,
			ChatID:        chat.ID,
			UserID:        user.ID,
		}
		err = suite.invitationsService.CancelInvitation(input)
		suite.ErrorIs(err, ErrInvitationNotExists)
	})
	suite.Run("приглашение могут отменить только инициатор, администратор и приглашаемый пользователь", func() {
		suite.Run("инициатор", func() {
			suite.invitationsService := newInvitationsService(t)

			chat := domain.Chat{
				ID: uuid.NewString(),
			}
			err := suite.invitationsService.ChatsRepo.Save(chat)
			suite.NoError(err)

			user := domain.User{
				ID: uuid.NewString(),
			}
			err = suite.invitationsService.UsersRepo.Save(user)
			suite.NoError(err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = suite.invitationsService.MembersRepo.Save(member)
			suite.NoError(err)

			invitation := domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member.UserID,
				UserID:        user.ID,
				ChatID:        chat.ID,
			}
			err = suite.invitationsService.InvitationsRepo.Save(invitation)
			suite.NoError(err)

			input := CancelInvitationInput{
				SubjectUserID: invitation.SubjectUserID,
				ChatID:        invitation.ChatID,
				UserID:        invitation.UserID,
			}
			err = suite.invitationsService.CancelInvitation(input)
			suite.NoError(err)
		})
		suite.Run("администратор", func() {
			suite.invitationsService := newInvitationsService(t)

			chatId := uuid.NewString()

			chiefMember := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chatId,
			}
			err := suite.invitationsService.MembersRepo.Save(chiefMember)
			suite.NoError(err)

			chat := domain.Chat{
				ID:          chatId,
				ChiefUserID: chiefMember.UserID,
			}
			err = suite.invitationsService.ChatsRepo.Save(chat)
			suite.NoError(err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = suite.invitationsService.MembersRepo.Save(member)
			suite.NoError(err)

			user := domain.User{
				ID: uuid.NewString(),
			}
			err = suite.invitationsService.UsersRepo.Save(user)
			suite.NoError(err)

			invitation := domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member.UserID,
				UserID:        user.ID,
				ChatID:        chat.ID,
			}
			err = suite.invitationsService.InvitationsRepo.Save(invitation)
			suite.NoError(err)

			input := CancelInvitationInput{
				SubjectUserID: chat.ChiefUserID,
				ChatID:        invitation.ChatID,
				UserID:        invitation.UserID,
			}
			err = suite.invitationsService.CancelInvitation(input)
			suite.NoError(err)
		})
		suite.Run("приглашаемый участник", func() {
			suite.invitationsService := newInvitationsService(t)

			chat := domain.Chat{
				ID: uuid.NewString(),
			}
			err := suite.invitationsService.ChatsRepo.Save(chat)
			suite.NoError(err)

			user := domain.User{
				ID: uuid.NewString(),
			}
			err = suite.invitationsService.UsersRepo.Save(user)
			suite.NoError(err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = suite.invitationsService.MembersRepo.Save(member)
			suite.NoError(err)

			invitation := domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member.UserID,
				UserID:        user.ID,
				ChatID:        chat.ID,
			}
			err = suite.invitationsService.InvitationsRepo.Save(invitation)
			suite.NoError(err)

			input := CancelInvitationInput{
				SubjectUserID: user.ID,
				ChatID:        chat.ID,
				UserID:        user.ID,
			}
			err = suite.invitationsService.CancelInvitation(input)
			suite.NoError(err)
		})
		suite.Run("посторонний участник чата", func() {
			suite.invitationsService := newInvitationsService(t)

			chat := domain.Chat{
				ID: uuid.NewString(),
			}
			err := suite.invitationsService.ChatsRepo.Save(chat)
			suite.NoError(err)

			user := domain.User{
				ID: uuid.NewString(),
			}
			err = suite.invitationsService.UsersRepo.Save(user)
			suite.NoError(err)

			member1 := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = suite.invitationsService.MembersRepo.Save(member1)
			suite.NoError(err)

			invitation := domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member1.UserID,
				UserID:        user.ID,
				ChatID:        chat.ID,
			}
			err = suite.invitationsService.InvitationsRepo.Save(invitation)
			suite.NoError(err)

			member2 := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = suite.invitationsService.MembersRepo.Save(member2)
			suite.NoError(err)

			input := CancelInvitationInput{
				SubjectUserID: member2.UserID,
				ChatID:        invitation.ChatID,
				UserID:        invitation.UserID,
			}
			err = suite.invitationsService.CancelInvitation(input)
			suite.ErrorIs(err, ErrSubjectUserNotAllowed)
		})
	})
	suite.Run("после отмены, в участник чата не добавляется", func() {

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(user)
		suite.NoError(err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = suite.invitationsService.MembersRepo.Save(member)
		suite.NoError(err)

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: member.UserID,
			UserID:        user.ID,
			ChatID:        chat.ID,
		}
		err = suite.invitationsService.InvitationsRepo.Save(invitation)
		suite.NoError(err)

		input := CancelInvitationInput{
			SubjectUserID: invitation.SubjectUserID,
			ChatID:        invitation.ChatID,
			UserID:        invitation.UserID,
		}
		err = suite.invitationsService.CancelInvitation(input)
		suite.NoError(err)

		members, err := suite.invitationsService.MembersRepo.List(domain.MembersFilter{})
		suite.NoError(err)
		if suite.Len(members, 1) {
			//assertEqualMembers(t, member, members[0])
		}
	})
	suite.Run("после отмены, приглашение удаляется", func() {

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := suite.invitationsService.ChatsRepo.Save(chat)
		suite.NoError(err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = suite.invitationsService.UsersRepo.Save(user)
		suite.NoError(err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = suite.invitationsService.MembersRepo.Save(member)
		suite.NoError(err)

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: member.UserID,
			UserID:        user.ID,
			ChatID:        chat.ID,
		}
		err = suite.invitationsService.InvitationsRepo.Save(invitation)
		suite.NoError(err)

		input := CancelInvitationInput{
			SubjectUserID: invitation.SubjectUserID,
			ChatID:        invitation.ChatID,
			UserID:        invitation.UserID,
		}
		err = suite.invitationsService.CancelInvitation(input)
		suite.NoError(err)

		invitations, err := suite.invitationsService.InvitationsRepo.List(domain.InvitationsFilter{})
		suite.NoError(err)
		suite.Len(invitations, 0)
	})
}
