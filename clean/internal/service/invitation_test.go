package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

func (suite *servicesTestSuite) saveInvitation(invitation domain.Invitation) domain.Invitation {
	err := suite.invitationsService.InvitationsRepo.Save(invitation)
	suite.Require().NoError(err)

	return invitation
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
	suite.Run("SubjectUserID является администратором", func() {
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

		suite.Run("созданные приглашения можно получить из репозитория", func() {
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
	})
	suite.Run("SubjectUserID не является администратором", func() {
		suite.Run("участник, не администратор, видит только им отправленные приглашения", func() {
			chat := suite.saveChat(domain.Chat{
				ID:          uuid.NewString(),
				ChiefUserID: uuid.NewString(),
			})

			member := suite.saveMember(domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			})

			invitationsSaved := make([]domain.Invitation, 3)
			for i := range invitationsSaved {
				invitationsSaved[i] = suite.saveInvitation(domain.Invitation{
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

			if suite.Len(invitationsFromService, len(invitationsSaved)) {
				for i, saved := range invitationsSaved {
					suite.Equal(saved, invitationsFromService[i])
				}
			}
		})
		suite.Run("если участника не существует", func() {
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
func Test_Invitations_UserInvitations(t *testing.T) {
	t.Run("пустой список из пустого репозитория", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		id := uuid.NewString()
		user := domain.User{
			ID: id,
		}
		err := serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)
		input := UserInvitationsInput{
			SubjectUserID: id,
			UserID:        id,
		}
		invs, err := serviceInvitations.UserInvitations(input)
		assert.NoError(t, err)
		assert.Len(t, invs, 0)
	})
	t.Run("пользователь должен существовать", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		id := uuid.NewString()
		input := UserInvitationsInput{
			SubjectUserID: id,
			UserID:        id,
		}
		invs, err := serviceInvitations.UserInvitations(input)
		assert.ErrorIs(t, err, ErrUserNotExists)
		assert.Len(t, invs, 0)
	})
	t.Run("пользователь может просматривать только свои приглашения", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		input := UserInvitationsInput{
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		invs, err := serviceInvitations.UserInvitations(input)
		assert.ErrorIs(t, err, ErrUnauthorizedInvitationsView)
		assert.Len(t, invs, 0)
	})
	t.Run("пустой список если у данного пользователя нету приглашений", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		for range 10 {
			inv := domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
			}
			err := serviceInvitations.InvitationsRepo.Save(inv)
			assert.NoError(t, err)
		}
		ourUserID := uuid.NewString()
		user := domain.User{
			ID: ourUserID,
		}
		err := serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)
		input := UserInvitationsInput{
			SubjectUserID: ourUserID,
			UserID:        ourUserID,
		}

		invs, err := serviceInvitations.UserInvitations(input)

		assert.NoError(t, err)
		assert.Len(t, invs, 0)
		allInvs, err := serviceInvitations.InvitationsRepo.List(domain.InvitationsFilter{})
		assert.Len(t, allInvs, 10)
		assert.NoError(t, err)
	})
	t.Run("у пользователя есть приглашение", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		userId := uuid.NewString()

		user := domain.User{
			ID: userId,
		}
		err := serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)

		input := UserInvitationsInput{
			SubjectUserID: userId,
			UserID:        userId,
		}
		chatId := uuid.NewString()
		err = serviceInvitations.InvitationsRepo.Save(domain.Invitation{
			ID:     uuid.NewString(),
			ChatID: chatId,
			UserID: userId,
		})
		assert.NoError(t, err)
		invs, err := serviceInvitations.UserInvitations(input)
		assert.NoError(t, err)
		if assert.Len(t, invs, 1) {
			assert.Equal(t, chatId, invs[0].ChatID)
			assert.Equal(t, userId, invs[0].UserID)
		}
	})
	t.Run("у пользователя несколько приглашений но не все из репозитория", func(t *testing.T) {
		const count = 5
		serviceInvitations := newInvitationsService(t)
		userId := uuid.NewString()
		user := domain.User{
			ID: userId,
		}
		err := serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)

		input := UserInvitationsInput{
			SubjectUserID: userId,
			UserID:        userId,
		}
		invsDomain := make([]domain.Invitation, count)
		for i := range count {
			inv := domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
				UserID: userId,
			}
			invsDomain[i] = inv
			err := serviceInvitations.InvitationsRepo.Save(invsDomain[i])
			assert.NoError(t, err)
		}
		for range count {
			err := serviceInvitations.InvitationsRepo.Save(domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: uuid.NewString(),
				UserID: uuid.NewString(),
			})
			assert.NoError(t, err)
		}

		invsRepo, err := serviceInvitations.UserInvitations(input)
		assert.NoError(t, err)
		if assert.Len(t, invsRepo, count) {
			for i, inv := range invsRepo {
				assert.Equal(t, inv.ID, invsDomain[i].ID)
				assert.Equal(t, inv.ChatID, invsDomain[i].ChatID)
				assert.Equal(t, inv.UserID, invsDomain[i].UserID)
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
	t.Run("участник отправляющий приглашения должен состоять в чате", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		subjectUser := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(subjectUser)
		assert.NoError(t, err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: subjectUser.ID,
		}
		err = serviceInvitations.MembersRepo.Save(member)
		assert.NoError(t, err)

		targetUser := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(targetUser)
		assert.NoError(t, err)

		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser.ID,
		}
		err = serviceInvitations.SendInvitation(input)
		assert.ErrorIs(t, err, ErrSubjectUserIsNotMember)
	})
	t.Run("UserID не должен состоять в чате ChatID", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = serviceInvitations.MembersRepo.Save(member)
		assert.NoError(t, err)

		targetUser := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(targetUser)
		assert.NoError(t, err)

		targetMember := domain.Member{
			ID:     uuid.NewString(),
			UserID: targetUser.ID,
			ChatID: chat.ID,
		}
		err = serviceInvitations.MembersRepo.Save(targetMember)
		assert.NoError(t, err)

		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser.ID,
		}
		err = serviceInvitations.SendInvitation(input)
		assert.ErrorIs(t, err, ErrUserAlreadyInChat)
	})
	t.Run("приглашать участников могут все члены чата", func(t *testing.T) {
		t.Run("админситратор", func(t *testing.T) {
			serviceInvitations := newInvitationsService(t)
			chatId := uuid.NewString()
			userId := uuid.NewString()

			chief := domain.Member{
				ID:     uuid.NewString(),
				UserID: userId,
				ChatID: chatId,
			}
			err := serviceInvitations.MembersRepo.Save(chief)
			assert.NoError(t, err)
			chat := domain.Chat{
				ID:          chatId,
				ChiefUserID: userId,
			}
			err = serviceInvitations.ChatsRepo.Save(chat)
			assert.NoError(t, err)

			targetUser := domain.User{
				ID: uuid.NewString(),
			}
			err = serviceInvitations.UsersRepo.Save(targetUser)
			assert.NoError(t, err)

			input := SendInvitationInput{
				ChatID:        chat.ID,
				SubjectUserID: chief.UserID,
				UserID:        targetUser.ID,
			}
			err = serviceInvitations.SendInvitation(input)
			assert.NoError(t, err)
		})
		t.Run("обычный участник чата", func(t *testing.T) {
			serviceInvitations := newInvitationsService(t)

			chat := domain.Chat{
				ID: uuid.NewString(),
			}
			err := serviceInvitations.ChatsRepo.Save(chat)
			assert.NoError(t, err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = serviceInvitations.MembersRepo.Save(member)
			assert.NoError(t, err)

			targetUser := domain.User{
				ID: uuid.NewString(),
			}
			err = serviceInvitations.UsersRepo.Save(targetUser)
			assert.NoError(t, err)

			input := SendInvitationInput{
				ChatID:        chat.ID,
				SubjectUserID: member.UserID,
				UserID:        targetUser.ID,
			}
			err = serviceInvitations.SendInvitation(input)
			assert.NoError(t, err)
		})
	})
	t.Run("UserID должен существовать", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = serviceInvitations.MembersRepo.Save(member)
		assert.NoError(t, err)

		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        uuid.NewString(),
		}
		err = serviceInvitations.SendInvitation(input)
		assert.ErrorIs(t, err, ErrUserNotExists)
	})
	t.Run("ChatID должен существовать", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		err := serviceInvitations.MembersRepo.Save(member)
		assert.NoError(t, err)

		targetUser := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(targetUser)
		assert.NoError(t, err)

		input := SendInvitationInput{
			ChatID:        uuid.NewString(),
			SubjectUserID: member.ID,
			UserID:        targetUser.ID,
		}
		err = serviceInvitations.SendInvitation(input)
		assert.ErrorIs(t, err, ErrChatNotExists)
	})
	t.Run("UserID нельзя приглашать более 1 раза", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = serviceInvitations.MembersRepo.Save(member)
		assert.NoError(t, err)

		targetUser := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(targetUser)
		assert.NoError(t, err)

		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser.ID,
		}
		err = serviceInvitations.SendInvitation(input)
		assert.NoError(t, err)

		err = serviceInvitations.SendInvitation(input)
		assert.ErrorIs(t, err, ErrUserAlreadyInviteInChat)
	})
	t.Run("можно приглашать больее 1 раза разных пользователей", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = serviceInvitations.MembersRepo.Save(member)
		assert.NoError(t, err)

		targetUser1 := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(targetUser1)
		assert.NoError(t, err)

		input1 := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser1.ID,
		}
		err = serviceInvitations.SendInvitation(input1)
		assert.NoError(t, err)

		targetUser2 := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(targetUser2)
		assert.NoError(t, err)

		input2 := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser2.ID,
		}

		err = serviceInvitations.SendInvitation(input2)
		assert.NoError(t, err)

		invsRepo, err := serviceInvitations.InvitationsRepo.List(domain.InvitationsFilter{})
		assert.NoError(t, err)
		assert.Len(t, invsRepo, 2)
		for i, invInput := range []SendInvitationInput{input1, input2} {
			assert.Equal(t, invInput.ChatID, invsRepo[i].ChatID)
			assert.Equal(t, invInput.SubjectUserID, invsRepo[i].SubjectUserID)
			assert.Equal(t, invInput.UserID, invsRepo[i].UserID)
		}
	})
}

func Test_AcceptInvitationInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		inp := AcceptInvitationInput{
			SubjectUserID: id,
			ChatID:        id,
		}
		return inp.Validate()
	})
}

func Test_Invitations_AcceptInvitation(t *testing.T) {
	t.Run("принятие не существующего приглашения", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)

		input := AcceptInvitationInput{
			SubjectUserID: user.ID,
			ChatID:        chat.ID,
		}
		err = serviceInvitations.AcceptInvitation(input)
		assert.ErrorIs(t, err, ErrInvitationNotExists)

		members, err := serviceInvitations.MembersRepo.List(domain.MembersFilter{})
		assert.NoError(t, err)
		assert.Len(t, members, 0)
	})
	t.Run("после принятия существующего приглашения, пользователь становится участником чата", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        user.ID,
			ChatID:        chat.ID,
		}
		err = serviceInvitations.InvitationsRepo.Save(invitation)
		assert.NoError(t, err)

		input := AcceptInvitationInput{
			SubjectUserID: user.ID,
			ChatID:        chat.ID,
		}
		err = serviceInvitations.AcceptInvitation(input)
		assert.NoError(t, err)

		members, err := serviceInvitations.MembersRepo.List(domain.MembersFilter{})
		assert.NoError(t, err)
		assert.Len(t, members, 1)
		assert.Equal(t, user.ID, members[0].UserID)
		assert.Equal(t, chat.ID, members[0].ChatID)
	})
	t.Run("принятие существующего приглашения в несуществющий чат", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)

		chatId := uuid.NewString()

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        user.ID,
			ChatID:        chatId,
		}
		err = serviceInvitations.InvitationsRepo.Save(invitation)
		assert.NoError(t, err)

		input := AcceptInvitationInput{
			SubjectUserID: user.ID,
			ChatID:        chatId,
		}
		err = serviceInvitations.AcceptInvitation(input)
		assert.ErrorIs(t, err, ErrChatNotExists)
	})
	t.Run("пользователя не существует", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
			ChatID:        chat.ID,
		}
		err = serviceInvitations.InvitationsRepo.Save(invitation)
		assert.NoError(t, err)

		input := AcceptInvitationInput{
			SubjectUserID: invitation.UserID,
			ChatID:        chat.ID,
		}
		err = serviceInvitations.AcceptInvitation(input)
		assert.ErrorIs(t, err, ErrUserNotExists)

		members, err := serviceInvitations.MembersRepo.List(domain.MembersFilter{})
		assert.NoError(t, err)
		assert.Len(t, members, 0)
	})
}

func Test_CancelInvitationInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := CancelInvitationInput{
			SubjectUserID: id,
			UserID:        id,
			ChatID:        id,
		}
		return input.Validate()
	})
}

func Test_Invitations_CancelInvitation(t *testing.T) {
	t.Run("отменить не существующее приглашение", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)

		input := CancelInvitationInput{
			SubjectUserID: user.ID,
			ChatID:        chat.ID,
			UserID:        user.ID,
		}
		err = serviceInvitations.CancelInvitation(input)
		assert.ErrorIs(t, err, ErrInvitationNotExists)
	})
	t.Run("приглашение могут отменить только инициатор, администратор и приглашаемый пользователь", func(t *testing.T) {
		t.Run("инициатор", func(t *testing.T) {
			serviceInvitations := newInvitationsService(t)

			chat := domain.Chat{
				ID: uuid.NewString(),
			}
			err := serviceInvitations.ChatsRepo.Save(chat)
			assert.NoError(t, err)

			user := domain.User{
				ID: uuid.NewString(),
			}
			err = serviceInvitations.UsersRepo.Save(user)
			assert.NoError(t, err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = serviceInvitations.MembersRepo.Save(member)
			assert.NoError(t, err)

			invitation := domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member.UserID,
				UserID:        user.ID,
				ChatID:        chat.ID,
			}
			err = serviceInvitations.InvitationsRepo.Save(invitation)
			assert.NoError(t, err)

			input := CancelInvitationInput{
				SubjectUserID: invitation.SubjectUserID,
				ChatID:        invitation.ChatID,
				UserID:        invitation.UserID,
			}
			err = serviceInvitations.CancelInvitation(input)
			assert.NoError(t, err)
		})
		t.Run("администратор", func(t *testing.T) {
			serviceInvitations := newInvitationsService(t)

			chatId := uuid.NewString()

			chiefMember := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chatId,
			}
			err := serviceInvitations.MembersRepo.Save(chiefMember)
			assert.NoError(t, err)

			chat := domain.Chat{
				ID:          chatId,
				ChiefUserID: chiefMember.UserID,
			}
			err = serviceInvitations.ChatsRepo.Save(chat)
			assert.NoError(t, err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = serviceInvitations.MembersRepo.Save(member)
			assert.NoError(t, err)

			user := domain.User{
				ID: uuid.NewString(),
			}
			err = serviceInvitations.UsersRepo.Save(user)
			assert.NoError(t, err)

			invitation := domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member.UserID,
				UserID:        user.ID,
				ChatID:        chat.ID,
			}
			err = serviceInvitations.InvitationsRepo.Save(invitation)
			assert.NoError(t, err)

			input := CancelInvitationInput{
				SubjectUserID: chat.ChiefUserID,
				ChatID:        invitation.ChatID,
				UserID:        invitation.UserID,
			}
			err = serviceInvitations.CancelInvitation(input)
			assert.NoError(t, err)
		})
		t.Run("приглашаемый участник", func(t *testing.T) {
			serviceInvitations := newInvitationsService(t)

			chat := domain.Chat{
				ID: uuid.NewString(),
			}
			err := serviceInvitations.ChatsRepo.Save(chat)
			assert.NoError(t, err)

			user := domain.User{
				ID: uuid.NewString(),
			}
			err = serviceInvitations.UsersRepo.Save(user)
			assert.NoError(t, err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = serviceInvitations.MembersRepo.Save(member)
			assert.NoError(t, err)

			invitation := domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member.UserID,
				UserID:        user.ID,
				ChatID:        chat.ID,
			}
			err = serviceInvitations.InvitationsRepo.Save(invitation)
			assert.NoError(t, err)

			input := CancelInvitationInput{
				SubjectUserID: user.ID,
				ChatID:        chat.ID,
				UserID:        user.ID,
			}
			err = serviceInvitations.CancelInvitation(input)
			assert.NoError(t, err)
		})
		t.Run("посторонний участник чата", func(t *testing.T) {
			serviceInvitations := newInvitationsService(t)

			chat := domain.Chat{
				ID: uuid.NewString(),
			}
			err := serviceInvitations.ChatsRepo.Save(chat)
			assert.NoError(t, err)

			user := domain.User{
				ID: uuid.NewString(),
			}
			err = serviceInvitations.UsersRepo.Save(user)
			assert.NoError(t, err)

			member1 := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = serviceInvitations.MembersRepo.Save(member1)
			assert.NoError(t, err)

			invitation := domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: member1.UserID,
				UserID:        user.ID,
				ChatID:        chat.ID,
			}
			err = serviceInvitations.InvitationsRepo.Save(invitation)
			assert.NoError(t, err)

			member2 := domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = serviceInvitations.MembersRepo.Save(member2)
			assert.NoError(t, err)

			input := CancelInvitationInput{
				SubjectUserID: member2.UserID,
				ChatID:        invitation.ChatID,
				UserID:        invitation.UserID,
			}
			err = serviceInvitations.CancelInvitation(input)
			assert.ErrorIs(t, err, ErrSubjectUserNotAllowed)
		})
	})
	t.Run("после отмены, в участник чата не добавляется", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = serviceInvitations.MembersRepo.Save(member)
		assert.NoError(t, err)

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: member.UserID,
			UserID:        user.ID,
			ChatID:        chat.ID,
		}
		err = serviceInvitations.InvitationsRepo.Save(invitation)
		assert.NoError(t, err)

		input := CancelInvitationInput{
			SubjectUserID: invitation.SubjectUserID,
			ChatID:        invitation.ChatID,
			UserID:        invitation.UserID,
		}
		err = serviceInvitations.CancelInvitation(input)
		assert.NoError(t, err)

		members, err := serviceInvitations.MembersRepo.List(domain.MembersFilter{})
		assert.NoError(t, err)
		if assert.Len(t, members, 1) {
			//assertEqualMembers(t, member, members[0])
		}
	})
	t.Run("после отмены, приглашение удаляется", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)

		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := serviceInvitations.ChatsRepo.Save(chat)
		assert.NoError(t, err)

		user := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(user)
		assert.NoError(t, err)

		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = serviceInvitations.MembersRepo.Save(member)
		assert.NoError(t, err)

		invitation := domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: member.UserID,
			UserID:        user.ID,
			ChatID:        chat.ID,
		}
		err = serviceInvitations.InvitationsRepo.Save(invitation)
		assert.NoError(t, err)

		input := CancelInvitationInput{
			SubjectUserID: invitation.SubjectUserID,
			ChatID:        invitation.ChatID,
			UserID:        invitation.UserID,
		}
		err = serviceInvitations.CancelInvitation(input)
		assert.NoError(t, err)

		invitations, err := serviceInvitations.InvitationsRepo.List(domain.InvitationsFilter{})
		assert.NoError(t, err)
		assert.Len(t, invitations, 0)
	})
}
