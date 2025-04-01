package service

import (
	"errors"
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

	return &Invitations{
		ChatsRepo:       chatsRepository,
		MembersRepo:     membersRepository,
		InvitationsRepo: invitationsRepository,
		History:         HistoryDummy{},
	}
}

func Test_ChatInvitationsInput_Validate(t *testing.T) {
	t.Run("UserID обязательное поле", func(t *testing.T) {
		input := ChatInvitationsInput{
			SubjectUserID: "",
			ChatID:        uuid.NewString(),
		}
		assert.Error(t, input.Validate())
	})
	t.Run("ChatID обязательное поле", func(t *testing.T) {
		input := ChatInvitationsInput{
			ChatID:        "",
			SubjectUserID: uuid.NewString(),
		}
		assert.Error(t, input.Validate())
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := ChatInvitationsInput{
			SubjectUserID: id,
			ChatID:        id,
		}
		return input.Validate()
	})
}

func Test_Invitations_ChatInvitations(t *testing.T) {
	t.Run("UserID должен быть администратором в чате с ChatID", func(t *testing.T) {
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
		invsChat, err := newService.ChatInvitations(input)
		assert.Error(t, err)
		assert.Len(t, invsChat, 0)
	})

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
		input := ChatInvitationsInput{
			SubjectUserID: userID,
			ChatID:        chatID,
		}

		const countInvs = 4
		exitsInvs := make([]domain.Invitation, 0, countInvs)
		err = nil
		for range countInvs {
			inv := domain.Invitation{
				ID:     uuid.NewString(),
				ChatID: chatID,
			}
			err = errors.Join(newService.InvitationsRepo.Save(inv))
			exitsInvs = append(exitsInvs, inv)
		}

		assert.NoError(t, err)
		invsChat, err := newService.ChatInvitations(input)
		assert.NoError(t, err)
		assert.Len(t, invsChat, countInvs)
		assert.Len(t, exitsInvs, countInvs)
		for i := range countInvs {
			assert.Equal(t, invsChat[i], exitsInvs[i])
		}
	})
}

// Test_UserInvitationsInput_Validate тестирует валидацию входящих параметров
func Test_UserInvitationsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := UserInvitationsInput{
			SubjectUserID: id,
			UserID:        uuid.NewString(),
		}
		return input.Validate()
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := UserInvitationsInput{
			SubjectUserID: uuid.NewString(),
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
		input := UserInvitationsInput{
			SubjectUserID: id,
			UserID:        id,
		}
		invs, err := serviceInvitations.UserInvitations(input)
		assert.NoError(t, err)
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
		invs, err := serviceInvitations.InvitationsRepo.List(domain.InvitationsFilter{
			ID: ourUserID,
		})
		assert.NoError(t, err)
		assert.Len(t, invs, 0)
		allInvs, err := serviceInvitations.InvitationsRepo.List(domain.InvitationsFilter{})
		assert.Len(t, allInvs, 10)
		assert.NoError(t, err)
	})
	t.Run("у пользователя есть приглашение", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		userId := uuid.NewString()
		input := UserInvitationsInput{
			SubjectUserID: userId,
			UserID:        userId,
		}
		chatId := uuid.NewString()
		err := serviceInvitations.InvitationsRepo.Save(domain.Invitation{
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

func Test_MemberSentInvitationsInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := MemberSentInvitationsInput{
			SubjectUserID: id,
			UserID:        uuid.NewString(),
			ChatID:        uuid.NewString(),
		}
		return in.Validate()
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := MemberSentInvitationsInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
			UserID:        id,
		}
		return in.Validate()
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := MemberSentInvitationsInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        id,
			UserID:        uuid.NewString(),
		}
		return in.Validate()
	})
}

func Test_Invitations_MemberSentInvitations(t *testing.T) {
	t.Run("получение списка для участника, у которого есть отправленные приглашения", func(t *testing.T) {
		const amount = 3
		serviceInvitations := newInvitationsService(t)
		chatId := uuid.NewString()
		err := serviceInvitations.ChatsRepo.Save(
			domain.Chat{
				ID: chatId,
			})
		assert.NoError(t, err)
		userId := uuid.NewString()
		input := MemberSentInvitationsInput{
			SubjectUserID: userId,
			UserID:        userId,
			ChatID:        chatId,
		}
		for range amount {
			inv := domain.Invitation{
				ID:     uuid.NewString(),
				UserID: userId,
				ChatID: chatId,
			}
			err := serviceInvitations.InvitationsRepo.Save(inv)
			assert.NoError(t, err)
		}
		invsRepo, err := serviceInvitations.MemberSentInvitations(input)
		assert.NoError(t, err)
		assert.Len(t, invsRepo, amount)
	})
	t.Run("UserID и SubjectUserID могут быть разными", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		chatId := uuid.NewString()
		err := serviceInvitations.ChatsRepo.Save(
			domain.Chat{
				ID: chatId,
			})
		assert.NoError(t, err)
		userId := uuid.NewString()
		input := MemberSentInvitationsInput{
			SubjectUserID: userId,
			UserID:        uuid.NewString(),
			ChatID:        chatId,
		}
		inv := domain.Invitation{
			ID:     uuid.NewString(),
			UserID: userId,
			ChatID: uuid.NewString(),
		}
		err = serviceInvitations.InvitationsRepo.Save(inv)
		assert.NoError(t, err)

		invsRepo, err := serviceInvitations.MemberSentInvitations(input)
		assert.NoError(t, err)
		assert.Len(t, invsRepo, 1)
	})
	t.Run("получение списка приглашений для участника без отправленных приглашений", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		userId := uuid.NewString()
		input := MemberSentInvitationsInput{
			SubjectUserID: userId,
			UserID:        userId,
		}
		inv := domain.Invitation{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		err := serviceInvitations.InvitationsRepo.Save(inv)
		assert.NoError(t, err)
		invsRepo, err := serviceInvitations.MemberSentInvitations(input)
		assert.NoError(t, err)
		assert.Len(t, invsRepo, 0)
	})
	t.Run("передача несуществующего UserID", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		input := MemberSentInvitationsInput{
			SubjectUserID: uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		inv := domain.Invitation{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		err := serviceInvitations.InvitationsRepo.Save(inv)
		assert.NoError(t, err)
		invsRepo, err := serviceInvitations.MemberSentInvitations(input)
		assert.NoError(t, err)
		assert.Len(t, invsRepo, 0)
	})
	t.Run("вызов метода с некоректными аргументами", func(t *testing.T) {
		serviceInvitations := newInvitationsService(t)
		input := MemberSentInvitationsInput{
			SubjectUserID: "",
			UserID:        uuid.NewString(),
		}
		inv := domain.Invitation{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		err := serviceInvitations.InvitationsRepo.Save(inv)
		assert.NoError(t, err)
		invsRepo, err := serviceInvitations.MemberSentInvitations(input)
		assert.Error(t, err)
		assert.Len(t, invsRepo, 0)
	})
}
