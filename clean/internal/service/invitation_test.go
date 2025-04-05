package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite/memory"
	"github.com/stretchr/testify/assert"
)

func newInvitationsService(t *testing.T) *Invitations {
	sqLiteInMemory, err := memory.Init(memory.Config{MigrationsDir: "../../migrations/repository/sqlite/memory"})
	assert.NoError(t, err)
	chatsRepository, err := sqLiteInMemory.NewChatsRepository()
	assert.NoError(t, err)
	membersRepository, err := sqLiteInMemory.NewMembersRepository()
	assert.NoError(t, err)
	invitationsRepository, err := sqLiteInMemory.NewInvitationsRepository()
	assert.NoError(t, err)
	usersRepository, err := sqLiteInMemory.NewUsersRepository()
	assert.NoError(t, err)

	return &Invitations{
		ChatsRepo:       chatsRepository,
		MembersRepo:     membersRepository,
		InvitationsRepo: invitationsRepository,
		UsersRepo:       usersRepository,
	}
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

func Test_Invitations_ChatInvitations(t *testing.T) {
	t.Run("SubjectUserID является администратором", func(t *testing.T) {
		t.Run("пустой список из чата без приглашений", func(t *testing.T) {
			newService := newInvitationsService(t)
			chatID := uuid.NewString()
			userID := uuid.NewString()
			err := newService.ChatsRepo.Save(domain.Chat{
				ID:          chatID,
				Name:        "Name1",
				ChiefUserID: userID,
			})
			assert.NoError(t, err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: userID,
				ChatID: chatID,
			}
			err = newService.MembersRepo.Save(member)
			assert.NoError(t, err)

			input := ChatInvitationsInput{
				SubjectUserID: userID,
				ChatID:        chatID,
			}
			invsChat, err := newService.ChatInvitations(input)
			assert.NoError(t, err)
			assert.Len(t, invsChat, 0)
		})

		t.Run("список из 4 приглашений из заполненного репозитория", func(t *testing.T) {
			newService := newInvitationsService(t)
			chatID := uuid.NewString()
			userID := uuid.NewString()
			err := newService.ChatsRepo.Save(domain.Chat{
				ID:          chatID,
				Name:        "Name1",
				ChiefUserID: userID,
			})
			assert.NoError(t, err)

			member := domain.Member{
				ID:     uuid.NewString(),
				UserID: userID,
				ChatID: chatID,
			}
			err = newService.MembersRepo.Save(member)
			assert.NoError(t, err)

			const countInvs = 4
			exitsInvs := make([]domain.Invitation, 0, countInvs)
			for range countInvs {
				inv := domain.Invitation{
					ID:     uuid.NewString(),
					ChatID: chatID,
				}
				err := newService.InvitationsRepo.Save(inv)
				assert.NoError(t, err)
				exitsInvs = append(exitsInvs, inv)
			}

			input := ChatInvitationsInput{
				SubjectUserID: userID,
				ChatID:        chatID,
			}

			invsChat, err := newService.ChatInvitations(input)
			assert.NoError(t, err)
			if assert.Len(t, invsChat, countInvs) {
				for i := range countInvs {
					assert.Equal(t, invsChat[i], exitsInvs[i])
				}
			}
		})
	})
	t.Run("SubjectUserID не является администратором", func(t *testing.T) {
		t.Run("если SubjectUserID не является администратором в чате с ChatID то должны возвращаться только пользователя приглашения",
			func(t *testing.T) {
				newService := newInvitationsService(t)
				chatID := uuid.NewString()

				err := newService.ChatsRepo.Save(domain.Chat{
					ID:          chatID,
					Name:        "Name1",
					ChiefUserID: uuid.NewString(),
				})
				assert.NoError(t, err)

				member := domain.Member{
					ID:     uuid.NewString(),
					UserID: uuid.NewString(),
					ChatID: chatID,
				}
				err = newService.MembersRepo.Save(member)
				assert.NoError(t, err)

				invitations := make([]domain.Invitation, 3)
				for i := range len(invitations) {
					invitation := domain.Invitation{
						ID:            uuid.NewString(),
						SubjectUserID: member.UserID,
						UserID:        uuid.NewString(),
						ChatID:        chatID,
					}
					invitations[i] = invitation
					err := newService.InvitationsRepo.Save(invitation)
					assert.NoError(t, err)
				}
				for range 3 {
					invitation := domain.Invitation{
						ID:            uuid.NewString(),
						SubjectUserID: uuid.NewString(),
						UserID:        uuid.NewString(),
						ChatID:        chatID,
					}
					err = newService.InvitationsRepo.Save(invitation)
					assert.NoError(t, err)
				}

				input := ChatInvitationsInput{
					ChatID:        chatID,
					SubjectUserID: member.UserID,
				}
				invsRepo, err := newService.ChatInvitations(input)
				assert.NoError(t, err)
				if assert.Len(t, invsRepo, len(invitations)) {
					for i := range len(invsRepo) {
						assert.Equal(t, invitations[i], invsRepo[i])
					}
				}
			})
		t.Run("если участника не существует", func(t *testing.T) {
			newService := newInvitationsService(t)
			chatID := uuid.NewString()

			err := newService.ChatsRepo.Save(domain.Chat{
				ID:          chatID,
				Name:        "Name1",
				ChiefUserID: uuid.NewString(),
			})
			assert.NoError(t, err)

			input := ChatInvitationsInput{
				ChatID:        chatID,
				SubjectUserID: uuid.NewString(),
			}
			invsRepo, err := newService.ChatInvitations(input)
			assert.Error(t, err)
			assert.Len(t, invsRepo, 0)
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
		input := SendChatInvitationInput{
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

		input := SendChatInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser.ID,
		}
		err = serviceInvitations.SendChatInvitation(input)
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

		input := SendChatInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser.ID,
		}
		err = serviceInvitations.SendChatInvitation(input)
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

			input := SendChatInvitationInput{
				ChatID:        chat.ID,
				SubjectUserID: chief.UserID,
				UserID:        targetUser.ID,
			}
			err = serviceInvitations.SendChatInvitation(input)
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

			input := SendChatInvitationInput{
				ChatID:        chat.ID,
				SubjectUserID: member.UserID,
				UserID:        targetUser.ID,
			}
			err = serviceInvitations.SendChatInvitation(input)
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

		input := SendChatInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        uuid.NewString(),
		}
		err = serviceInvitations.SendChatInvitation(input)
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

		input := SendChatInvitationInput{
			ChatID:        uuid.NewString(),
			SubjectUserID: member.ID,
			UserID:        targetUser.ID,
		}
		err = serviceInvitations.SendChatInvitation(input)
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

		input := SendChatInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser.ID,
		}
		err = serviceInvitations.SendChatInvitation(input)
		assert.NoError(t, err)

		err = serviceInvitations.SendChatInvitation(input)
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

		input1 := SendChatInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser1.ID,
		}
		err = serviceInvitations.SendChatInvitation(input1)
		assert.NoError(t, err)

		targetUser2 := domain.User{
			ID: uuid.NewString(),
		}
		err = serviceInvitations.UsersRepo.Save(targetUser2)
		assert.NoError(t, err)

		input2 := SendChatInvitationInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
			UserID:        targetUser2.ID,
		}

		err = serviceInvitations.SendChatInvitation(input2)
		assert.NoError(t, err)

		invsRepo, err := serviceInvitations.InvitationsRepo.List(domain.InvitationsFilter{})
		assert.NoError(t, err)
		assert.Len(t, invsRepo, 2)
		for i, invInput := range []SendChatInvitationInput{input1, input2} {
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
			assert.ErrorIs(t, err, ErrSubjectUserIsNotChief)
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
			assertEqualMembers(t, member, members[0])
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
