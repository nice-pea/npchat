package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite/memory"
)

// newChatsService создает объект сервиса Chats с sqlite/memory репозиториями
func newMembersService(t *testing.T) *Members {
	sqLiteInMemory, err := memory.Init(memory.Config{MigrationsDir: "../../migrations/repository/sqlite/memory"})
	assert.NoError(t, err)
	chatsRepository, err := sqLiteInMemory.NewChatsRepository()
	assert.NoError(t, err)
	membersRepository, err := sqLiteInMemory.NewMembersRepository()
	assert.NoError(t, err)
	return &Members{
		ChatsRepo:   chatsRepository,
		MembersRepo: membersRepository,
	}
}

func assertEqualMembers(t *testing.T, expected, actual domain.Member) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.ChatID, actual.ChatID)
}

func Test_ChatMembersInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := ChatMembersInput{
			ChatID:        id,
			SubjectUserID: uuid.NewString(),
		}
		return in.Validate()
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := ChatMembersInput{
			ChatID:        uuid.NewString(),
			SubjectUserID: id,
		}
		return in.Validate()
	})
}

func Test_Members_ChatMembers(t *testing.T) {
	t.Run("чат должен существовать", func(t *testing.T) {
		membersService := newMembersService(t)
		input := ChatMembersInput{
			ChatID:        uuid.New().String(),
			SubjectUserID: uuid.NewString(),
		}
		members, err := membersService.ChatMembers(input)
		assert.Error(t, err)
		assert.Len(t, members, 0)
	})
	t.Run("пользователь должен быть участником чата", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{ID: uuid.NewString()}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: uuid.NewString(),
		}
		err = membersService.MembersRepo.Save(member)
		assert.NoError(t, err)
		input := ChatMembersInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
		}
		members, err := membersService.ChatMembers(input)
		assert.Error(t, err)
		assert.Len(t, members, 0)
	})
	t.Run("возвращается список участников чата", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{ID: uuid.NewString()}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		members := make([]domain.Member, 20)
		for i := range members {
			members[i] = domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			}
			err = membersService.MembersRepo.Save(members[i])
			assert.NoError(t, err)
		}
		subjectMember := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(subjectMember)
		assert.NoError(t, err)
		input := ChatMembersInput{
			ChatID:        subjectMember.ChatID,
			SubjectUserID: subjectMember.UserID,
		}
		membersFromRepo, err := membersService.ChatMembers(input)
		assert.NoError(t, err)
		if assert.Len(t, membersFromRepo, len(members)+1) {
			for i := range members {
				assertEqualMembers(t, members[i], membersFromRepo[i])
			}
			lastRepoMember := membersFromRepo[len(membersFromRepo)-1]
			assertEqualMembers(t, subjectMember, lastRepoMember)
		}
	})
}

/*
Покинуть чат

	пользователь должен быть участникам чата
	пользователь должен не должен быть главным администратором
	Входящие параметры должны валидироваться

Принудительно удалить участника

	пользователь должен быть участникам чата
	пользователь должен быть главным администратором
	нельзя удалить самого себя
	Входящие параметры должны валидироваться
*/

func Test_LeaveInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := LeaveInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        id,
		}
		return in.Validate()
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := LeaveInput{
			SubjectUserID: id,
			ChatID:        uuid.NewString(),
		}
		return in.Validate()
	})
}

func Test_Members_Leave(t *testing.T) {
	t.Run("чат должен существовать", func(t *testing.T) {
		membersService := newMembersService(t)
		input := LeaveInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
		}
		err := membersService.Leave(input)
		assert.Error(t, err)
	})
	t.Run("пользователь должен быть участником чата", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		input := LeaveInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        chat.ID,
		}
		err = membersService.Leave(input)
		assert.Error(t, err)
	})
	t.Run("пользователь не должен быть главным администратором чата", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{ID: uuid.NewString(), ChiefUserID: uuid.New().String()}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(member)
		assert.NoError(t, err)
		input := LeaveInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		err = membersService.Leave(input)
		assert.Error(t, err)
	})
	t.Run("после выхода пользователь перестает быть участником", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{ID: uuid.NewString(), ChiefUserID: uuid.NewString()}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(member)
		assert.NoError(t, err)
		input := LeaveInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		err = membersService.Leave(input)
		assert.NoError(t, err)
		membersFilter := domain.MembersFilter{ID: member.ID}
		members, err := membersService.MembersRepo.List(membersFilter)
		assert.NoError(t, err)
		assert.Len(t, members, 0)
	})
}

func Test_Members_Delete(t *testing.T) {
	_ = []string{
		"чат должен существовать",
		"пользователь должен быть участником чата",
		"пользователь должен быть главным администратором чата",
		"удаляемый участник должен существовать",
		"",
	}
}
