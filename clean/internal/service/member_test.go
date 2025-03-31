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

// assertEqualMembers сравнивает поля domain.Member
func assertEqualMembers(t *testing.T, expected, actual domain.Member) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.ChatID, actual.ChatID)
}

// Test_ChatMembersInput_Validate тестирует валидацию входящих параметров
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

// Test_Members_ChatMembers тестирует получение списка участников чата
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

// Test_LeaveInput_Validate тестирует валидацию входящих параметров
func Test_LeaveInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := LeaveChatInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        id,
		}
		return in.Validate()
	})
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := LeaveChatInput{
			SubjectUserID: id,
			ChatID:        uuid.NewString(),
		}
		return in.Validate()
	})
}

// Test_Members_LeaveChat тестирует выход участника из чата
func Test_Members_LeaveChat(t *testing.T) {
	t.Run("чат должен существовать", func(t *testing.T) {
		membersService := newMembersService(t)
		input := LeaveChatInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
		}
		err := membersService.LeaveChat(input)
		assert.Error(t, err)
	})
	t.Run("пользователь должен быть участником чата", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{
			ID: uuid.NewString(),
		}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		input := LeaveChatInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        chat.ID,
		}
		err = membersService.LeaveChat(input)
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
		input := LeaveChatInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		err = membersService.LeaveChat(input)
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
		input := LeaveChatInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		err = membersService.LeaveChat(input)
		assert.NoError(t, err)
		membersFilter := domain.MembersFilter{ID: member.ID}
		members, err := membersService.MembersRepo.List(membersFilter)
		assert.NoError(t, err)
		assert.Len(t, members, 0)
	})
}

// Test_DeleteMemberInput_Validate тестирует валидацию входящих параметров
func Test_DeleteMemberInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := DeleteMemberInput{
			SubjectUserID: id,
			ChatID:        id,
			UserID:        id,
		}
		return input.Validate()
	})
}

// Test_Members_DeleteMember тестирует удаление участника чата
func Test_Members_DeleteMember(t *testing.T) {
	t.Run("чат должен существовать", func(t *testing.T) {
		membersService := newMembersService(t)
		input := DeleteMemberInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		err := membersService.DeleteMember(input)
		assert.Error(t, err)
	})
	t.Run("удаляемый участник должен существовать", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{ID: uuid.NewString()}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(member)
		input := DeleteMemberInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
			UserID:        uuid.NewString(),
		}
		err = membersService.DeleteMember(input)
		assert.Error(t, err)
	})
	t.Run("пользователь должен быть участником чата", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{ID: uuid.NewString()}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(member)
		input := DeleteMemberInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        chat.ID,
			UserID:        member.UserID,
		}
		err = membersService.DeleteMember(input)
		assert.Error(t, err)
	})
	t.Run("нельзя удалить самого себя", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{ID: uuid.NewString()}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		member := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(member)
		assert.NoError(t, err)
		input := DeleteMemberInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
			UserID:        member.UserID,
		}
		err = membersService.DeleteMember(input)
		assert.Error(t, err)
	})
	t.Run("пользователь должен быть главным администратором чата", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{ID: uuid.NewString()}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		memberForDelete := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(memberForDelete)
		subjectMember := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(subjectMember)
		input := DeleteMemberInput{
			SubjectUserID: subjectMember.UserID,
			ChatID:        chat.ID,
			UserID:        memberForDelete.UserID,
		}
		err = membersService.DeleteMember(input)
		assert.Error(t, err)
	})
	t.Run("после удаления участник перестает быть участником", func(t *testing.T) {
		membersService := newMembersService(t)
		chat := domain.Chat{ID: uuid.NewString(), ChiefUserID: uuid.NewString()}
		err := membersService.ChatsRepo.Save(chat)
		assert.NoError(t, err)
		memberForDelete := domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(memberForDelete)
		subjectMember := domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		}
		err = membersService.MembersRepo.Save(subjectMember)
		input := DeleteMemberInput{
			SubjectUserID: subjectMember.UserID,
			ChatID:        chat.ID,
			UserID:        memberForDelete.UserID,
		}
		err = membersService.DeleteMember(input)
		assert.NoError(t, err)
		membersFilter := domain.MembersFilter{ID: memberForDelete.ID}
		members, err := membersService.MembersRepo.List(membersFilter)
		assert.NoError(t, err)
		assert.Len(t, members, 0)
	})
}
