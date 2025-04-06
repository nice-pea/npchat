package service

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite/memory"
)

type membersTestEnv struct {
	membersService *Members
	t              *testing.T
}

func initMembersTestEnv(t *testing.T) membersTestEnv {
	env := membersTestEnv{
		membersService: &Members{},
		t:              t,
	}
	sqLiteInMemory, err := memory.Init(memory.Config{MigrationsDir: "../../migrations/repository/sqlite/memory"})
	assert.NoError(t, err)
	env.membersService.ChatsRepo, err = sqLiteInMemory.NewChatsRepository()
	assert.NoError(t, err)
	env.membersService.MembersRepo, err = sqLiteInMemory.NewMembersRepository()
	assert.NoError(t, err)

	return env
}

// assertEqualMembers сравнивает поля domain.Member
func (e *membersTestEnv) assertEqualMembers(expected, actual domain.Member) {
	assert.Equal(e.t, expected.ID, actual.ID)
	assert.Equal(e.t, expected.UserID, actual.UserID)
	assert.Equal(e.t, expected.ChatID, actual.ChatID)
}

func (e *membersTestEnv) saveChat(chat domain.Chat) domain.Chat {
	err := e.membersService.ChatsRepo.Save(chat)
	assert.NoError(e.t, err)

	return chat
}

func (e *membersTestEnv) saveMember(member domain.Member) domain.Member {
	err := e.membersService.MembersRepo.Save(member)
	assert.NoError(e.t, err)

	return member
}

// Test_ChatMembersInput_Validate тестирует валидацию входящих параметров
func Test_ChatMembersInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := ChatMembersInput{
			ChatID:        id,
			SubjectUserID: id,
		}
		return in.Validate()
	})
}

// Test_Members_ChatMembers тестирует получение списка участников чата
func Test_Members_ChatMembers(t *testing.T) {
	t.Run("чат должен существовать", func(t *testing.T) {
		env := initMembersTestEnv(t)
		input := ChatMembersInput{
			ChatID:        uuid.NewString(),
			SubjectUserID: uuid.NewString(),
		}
		// Запросить список участников чата
		members, err := env.membersService.ChatMembers(input)
		assert.ErrorIs(t, err, ErrChatNotExists)
		assert.Empty(t, members)
	})
	t.Run("пользователь должен быть участником чата", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Создать чат
		chat := env.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать участника в другом чате
		member := env.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: uuid.NewString(),
		})
		// Запросить список участников чата
		input := ChatMembersInput{
			ChatID:        chat.ID,
			SubjectUserID: member.UserID,
		}
		members, err := env.membersService.ChatMembers(input)
		// Вернется ошибка, потому пользователь не является участником чата
		assert.ErrorIs(t, err, ErrSubjectUserIsNotMember)
		assert.Empty(t, members)
	})
	t.Run("возвращается список участников чата", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Создать чат
		chat := env.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать несколько участников в чате
		const membersAllCount = 20
		savedMembers := make([]domain.Member, membersAllCount)
		for i := range membersAllCount {
			// Создать участника в чате
			savedMembers[i] = env.saveMember(domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			})
		}
		// Запрашивать список будет первый участник
		subjectMember := savedMembers[0]
		// Получить список участников в чате
		input := ChatMembersInput{
			ChatID:        subjectMember.ChatID,
			SubjectUserID: subjectMember.UserID,
		}
		membersFromRepo, err := env.membersService.ChatMembers(input)
		assert.NoError(t, err)
		if assert.Len(t, membersFromRepo, membersAllCount) {
			// Сравнить каждого сохраненного участника с ранее созданным
			for i := range savedMembers {
				env.assertEqualMembers(savedMembers[i], membersFromRepo[i])
			}
		}
	})
}

// Test_LeaveInput_Validate тестирует валидацию входящих параметров
func Test_LeaveInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := LeaveChatInput{
			SubjectUserID: id,
			ChatID:        id,
		}
		return in.Validate()
	})
}

// Test_Members_LeaveChat тестирует выход участника из чата
func Test_Members_LeaveChat(t *testing.T) {
	t.Run("чат должен существовать", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Покинуть чат
		input := LeaveChatInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
		}
		err := env.membersService.LeaveChat(input)
		// Вернется ошибка, потому что чата не существует
		assert.ErrorIs(t, err, ErrChatNotExists)
	})
	t.Run("пользователь должен быть участником чата", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Создать чат
		chat := env.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Покинуть чат
		input := LeaveChatInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        chat.ID,
		}
		err := env.membersService.LeaveChat(input)
		// Вернется ошибка, потому что пользователь не участник чата
		assert.ErrorIs(t, err, ErrSubjectUserIsNotMember)
	})
	t.Run("пользователь не должен быть главным администратором чата", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Создать чат
		chat := env.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Создать участника главного администратора в этом чате
		member := env.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		})
		// Покинуть чат
		input := LeaveChatInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		err := env.membersService.LeaveChat(input)
		// Вернется ошибка, потому что пользователь главный администратор чата
		assert.ErrorIs(t, err, ErrSubjectUserShouldNotBeChief)
	})
	t.Run("после выхода пользователь перестает быть участником", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Создать чат
		chat := env.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать участника в этом чате
		member := env.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Покинуть чат
		input := LeaveChatInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
		}
		err := env.membersService.LeaveChat(input)
		assert.NoError(t, err)
		// Получить список участников чата
		membersFilter := domain.MembersFilter{ID: member.ID}
		members, err := env.membersService.MembersRepo.List(membersFilter)
		assert.NoError(t, err)
		// В чате не осталось участников
		assert.Empty(t, members)
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
	t.Run("нельзя удалить самого себя", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Удалить участника
		userID := uuid.NewString()
		input := DeleteMemberInput{
			SubjectUserID: userID,
			ChatID:        uuid.NewString(),
			UserID:        userID,
		}
		err := env.membersService.DeleteMember(input)
		// Вернется ошибка, потому что пользователь пытается удалить самого себя
		assert.ErrorIs(t, err, ErrMemberCannotDeleteHimself)
	})
	t.Run("чат должен существовать", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Удалить участника
		input := DeleteMemberInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        uuid.NewString(),
			UserID:        uuid.NewString(),
		}
		err := env.membersService.DeleteMember(input)
		// Вернется ошибка, потому что чата не существует
		assert.ErrorIs(t, err, ErrChatNotExists)
	})
	t.Run("subject должен быть участником чата", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Создать чат
		chat := env.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Удалить участника
		input := DeleteMemberInput{
			SubjectUserID: uuid.NewString(),
			ChatID:        chat.ID,
			UserID:        uuid.NewString(),
		}
		err := env.membersService.DeleteMember(input)
		// Вернется ошибка, потому что пользователь не участник чата
		assert.ErrorIs(t, err, ErrSubjectUserIsNotMember)
	})
	t.Run("subject должен быть главным администратором чата", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Создать чат
		chat := env.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать участника
		subjectMember := env.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Удалить участника
		input := DeleteMemberInput{
			SubjectUserID: subjectMember.UserID,
			ChatID:        chat.ID,
			UserID:        uuid.NewString(),
		}
		err := env.membersService.DeleteMember(input)
		// Вернется ошибка, потому что участник не главный администратор
		assert.ErrorIs(t, err, ErrSubjectUserIsNotChief)
	})
	t.Run("user должен быть участником чата", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Создать чат
		chat := env.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Создать участника
		member := env.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		})
		// Удалить участника
		input := DeleteMemberInput{
			SubjectUserID: member.UserID,
			ChatID:        chat.ID,
			UserID:        uuid.NewString(),
		}
		err := env.membersService.DeleteMember(input)
		// Вернется ошибка, потому что удаляемый пользователь не является участником
		assert.ErrorIs(t, err, ErrUserIsNotMember)
	})
	t.Run("после удаления участник перестает быть участником", func(t *testing.T) {
		env := initMembersTestEnv(t)
		// Создать чат
		chat := env.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Создать участника для удаления
		memberForDelete := env.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Создать участника
		subjectMember := env.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		})
		// Удалить участника
		input := DeleteMemberInput{
			SubjectUserID: subjectMember.UserID,
			ChatID:        chat.ID,
			UserID:        memberForDelete.UserID,
		}
		err := env.membersService.DeleteMember(input)
		assert.NoError(t, err)
		// Найти удаленного участника
		membersFilter := domain.MembersFilter{ID: memberForDelete.ID}
		members, err := env.membersService.MembersRepo.List(membersFilter)
		assert.NoError(t, err)
		// Такого участника больше нет
		assert.Empty(t, members)
	})
}
