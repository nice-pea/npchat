package service

import (
	"testing"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

// saveInvitation сохраняет приглашение в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveInvitation(invitation domain.Invitation) domain.Invitation {
	err := suite.invitationsService.InvitationsRepo.Save(invitation)
	suite.Require().NoError(err)

	return invitation
}

// saveUser сохраняет пользователя в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveUser(user domain.User) domain.User {
	err := suite.invitationsService.UsersRepo.Save(user)
	suite.Require().NoError(err)

	return user
}

// Test_ChatInvitationsInput_Validate тестирует валидацию входящих параметров
func Test_ChatInvitationsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := ChatInvitationsInput{
			SubjectUserID: id,
			ChatID:        id,
		}
		return input.Validate()
	})
}

// Test_Invitations_ChatInvitations тестирует получение списка приглашений
func (suite *servicesTestSuite) Test_Invitations_ChatInvitations() {
	suite.Run("чат должен существовать", func() {
		// Получить список приглашений
		input := ChatInvitationsInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
		}
		invitations, err := suite.invitationsService.ChatInvitations(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
		suite.Empty(invitations)
	})

	suite.Run("субъект должен быть участником чата", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Получить список приглашений
		input := ChatInvitationsInput{
			ChatID:        chat.ID,
			SubjectUserID: uuid.NewString(),
		}
		invitations, err := suite.invitationsService.ChatInvitations(input)
		// Вернется ошибка, потому что пользователь не участник чата
		suite.ErrorIs(err, ErrSubjectUserIsNotMember)
		suite.Empty(invitations)
	})

	suite.Run("пустой список из чата без приглашений", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать чат
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Получить список приглашений
		input := ChatInvitationsInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		invitations, err := suite.invitationsService.ChatInvitations(input)
		suite.NoError(err)
		suite.Empty(invitations)
	})

	suite.Run("субъект не администратор чата и видит только отправленные им приглашения", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Создать участника в чате
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
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
		input := ChatInvitationsInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
		}
		invitationsFromService, err := suite.invitationsService.ChatInvitations(input)
		suite.Require().NoError(err)
		// В списке будут приглашения отправленные участником
		if suite.Len(invitationsFromService, len(subjectInvitations)) {
			for i, subjectInvitation := range subjectInvitations {
				suite.Equal(subjectInvitation, invitationsFromService[i])
			}
		}
	})

	suite.Run("субъект является администратором чата и видит все отправленные приглашения в чат", func() {
		// Создать чат с указанием главного администратора
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Создать участника для главного администратора в чате
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		})
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
		input := ChatInvitationsInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		invitationFromService, err := suite.invitationsService.ChatInvitations(input)
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
		input := UserInvitationsInput{
			SubjectUserID: id,
			UserID:        id,
		}
		return input.Validate()
	})
}

// Test_Invitations_UserInvitations тестирует получение списка приглашений направленных пользователю
func (suite *servicesTestSuite) Test_Invitations_UserInvitations() {
	suite.Run("пользователь должен существовать", func() {
		id := uuid.NewString()
		// Получить список приглашений
		input := UserInvitationsInput{
			SubjectUserID: id,
			UserID:        id,
		}
		invitations, err := suite.invitationsService.UserInvitations(input)
		// Вернется ошибка, потому что пользователя не существует
		suite.ErrorIs(err, ErrUserNotExists)
		suite.Empty(invitations)
	})

	suite.Run("пользователь может просматривать только свои приглашения", func() {
		// Получить список приглашений
		input := UserInvitationsInput{
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		invitations, err := suite.invitationsService.UserInvitations(input)
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
		input := UserInvitationsInput{
			SubjectUserID: user.ID,
			UserID:        user.ID,
		}
		invitations, err := suite.invitationsService.UserInvitations(input)
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
		input := UserInvitationsInput{
			SubjectUserID: user.ID,
			UserID:        user.ID,
		}
		invitationsFromRepo, err := suite.invitationsService.UserInvitations(input)
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
		input := SendInvitationInput{
			SubjectUserID: id,
			ChatID:        id,
			UserID:        id,
		}
		return input.Validate()
	})
}

// Test_Invitations_SendChatInvitation тестирует отправку приглашения
func (suite *servicesTestSuite) Test_Invitations_SendChatInvitation() {
	suite.Run("чат должен существовать", func() {
		// Отправить приглашение
		input := SendInvitationInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		invitation, err := suite.invitationsService.SendInvitation(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
		suite.Zero(invitation)
	})

	suite.Run("субъект должен быть участником", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Отправить приглашение
		input := SendInvitationInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        chat.ID,
			UserID:        uuid.NewString(),
		}
		invitation, err := suite.invitationsService.SendInvitation(input)
		// Вернется ошибка, потому что субъект не является участником чата
		suite.ErrorIs(err, ErrSubjectUserIsNotMember)
		suite.Zero(invitation)
	})

	suite.Run("приглашаемый пользователь должен существовать", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать участника
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Отправить приглашение
		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        uuid.NewString(),
		}
		invitation, err := suite.invitationsService.SendInvitation(input)
		// Вернется ошибка, потому что приглашаемого пользователя не существует
		suite.ErrorIs(err, ErrUserNotExists)
		suite.Zero(invitation)
	})

	suite.Run("приглашаемый пользователь не должен состоять в этом чате", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
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
		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: subjectMember.UserID,
			UserID:        targetUser.ID,
		}
		invitation, err := suite.invitationsService.SendInvitation(input)
		// Вернется ошибка, потому что приглашаемый пользователь уже является участником этого чата
		suite.ErrorIs(err, ErrUserIsAlreadyInChat)
		suite.Zero(invitation)
	})

	suite.Run("одновременно не может существовать несколько приглашений одного пользователя в этот чат", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
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
		input := SendInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: subjectMember.UserID,
			UserID:        targetUser.ID,
		}
		invitation, err := suite.invitationsService.SendInvitation(input)
		suite.NoError(err)
		suite.Require().NotZero(invitation)
		// Отправить повторно приглашение
		invitation, err = suite.invitationsService.SendInvitation(input)
		// Вернется ошибка, потому что этот пользователь уже приглашен в чат
		suite.ErrorIs(err, ErrUserIsAlreadyInvited)
		suite.Zero(invitation)
	})

	suite.Run("любой участник может приглашать много пользователей", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
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
				input := SendInvitationInput{
					ChatID:        chat.ID,
					SubjectUserID: subjectMember.UserID,
					UserID:        targetUser.ID,
				}
				invitation, err := suite.invitationsService.SendInvitation(input)
				suite.NoError(err)
				suite.Require().NotZero(invitation)
				createdInvitations = append(createdInvitations, invitation)
			}
		}
		// Получить список приглашений
		invitationsFromRepo, err := suite.invitationsService.InvitationsRepo.List(domain.InvitationsFilter{})
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
		inp := AcceptInvitationInput{
			SubjectUserID: id,
			ChatID:        id,
		}
		return inp.Validate()
	})
}

// Test_Invitations_AcceptInvitation тестирует принятие приглашения
func (suite *servicesTestSuite) Test_Invitations_AcceptInvitation() {
	suite.Run("чат должен существовать", func() {
		// Принять приглашение
		input := AcceptInvitationInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
		}
		err := suite.invitationsService.AcceptInvitation(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
	})

	suite.Run("приглашение должно существовать", func() {
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Принять приглашение
		input := AcceptInvitationInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        chat.ID,
		}
		err := suite.invitationsService.AcceptInvitation(input)
		suite.ErrorIs(err, ErrInvitationNotExists)
	})

	suite.Run("пользователь должен существовать", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать приглашение
		invitation := suite.saveInvitation(domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
			ChatID:        chat.ID,
		})
		// Принять приглашение
		input := AcceptInvitationInput{
			SubjectUserID: invitation.UserID,
			ChatID:        chat.ID,
		}
		err := suite.invitationsService.AcceptInvitation(input)
		suite.ErrorIs(err, ErrUserNotExists)
	})

	suite.Run("приняв приглашение, пользователь становится участником чата", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		user := suite.saveUser(domain.User{
			ID: uuid.NewString(),
		})
		// Создать приглашение
		suite.saveInvitation(domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        user.ID,
			ChatID:        chat.ID,
		})
		// Принять приглашение
		input := AcceptInvitationInput{
			SubjectUserID: user.ID,
			ChatID:        chat.ID,
		}
		err := suite.invitationsService.AcceptInvitation(input)
		suite.NoError(err)
		// Получить список участников
		members, err := suite.invitationsService.MembersRepo.List(domain.MembersFilter{})
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
		input := CancelInvitationInput{
			SubjectUserID: id,
			UserID:        id,
			ChatID:        id,
		}
		return input.Validate()
	})
}

// Test_Invitations_CancelInvitation тестирует отмену приглашения
func (suite *servicesTestSuite) Test_Invitations_CancelInvitation() {
	suite.Run("чат должен существовать", func() {
		// Отменить приглашение
		input := CancelInvitationInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		err := suite.invitationsService.CancelInvitation(input)
		// Вернется ошибка потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
	})

	suite.Run("приглашение должно существовать", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Отменить приглашение
		input := CancelInvitationInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        chat.ID,
			UserID:        uuid.NewString(),
		}
		err := suite.invitationsService.CancelInvitation(input)
		// Вернется ошибка потому что приглашения не существует
		suite.ErrorIs(err, ErrInvitationNotExists)
	})

	suite.Run("приглашение могут отменить только пригласивший и приглашаемый пользователи, и администратор чата", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
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
			suite.saveInvitation(domain.Invitation{
				ID:            uuid.NewString(),
				SubjectUserID: subjectMember.UserID,
				UserID:        userID,
				ChatID:        chat.ID,
			})
			// Отменить приглашение
			input := CancelInvitationInput{
				SubjectUserID: subjectUserID,
				ChatID:        chat.ID,
				UserID:        userID,
			}
			err := suite.invitationsService.CancelInvitation(input)
			suite.Require().NoError(err)
		}
	})

	suite.Run("другие участники не могут отменять приглашать ", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		userID := uuid.NewString()
		// Случайный участника
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
		})
		// Создать приглашение
		suite.saveInvitation(domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: uuid.NewString(),
			UserID:        userID,
			ChatID:        chat.ID,
		})
		// Отменить приглашение
		input := CancelInvitationInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
			UserID:        userID,
		}
		err := suite.invitationsService.CancelInvitation(input)
		// Вернется ошибка, потому что случайный участник не может отменять приглашение
		suite.ErrorIs(err, ErrSubjectUserNotAllowed)
	})

	suite.Run("после отмены, приглашение удаляется", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать участника
		subjectMember := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Создать приглашение
		invitation := suite.saveInvitation(domain.Invitation{
			ID:            uuid.NewString(),
			SubjectUserID: subjectMember.UserID,
			UserID:        uuid.NewString(),
			ChatID:        chat.ID,
		})
		// Отменить приглашение
		input := CancelInvitationInput{
			SubjectUserID: invitation.SubjectUserID,
			ChatID:        invitation.ChatID,
			UserID:        invitation.UserID,
		}
		err := suite.invitationsService.CancelInvitation(input)
		suite.Require().NoError(err)
		// Получить список приглашений
		invitations, err := suite.invitationsService.InvitationsRepo.List(domain.InvitationsFilter{})
		suite.NoError(err)
		suite.Empty(invitations)
	})
}
