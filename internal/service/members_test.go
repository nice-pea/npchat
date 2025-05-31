package service

import (
	"testing"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/helpers_tests"
)

// assertEqualMembers сравнивает поля domain.Member
func (suite *servicesTestSuite) assertEqualMembers(expected, actual domain.Member) {
	suite.Equal(expected.ID, actual.ID)
	suite.Equal(expected.UserID, actual.UserID)
	suite.Equal(expected.ChatID, actual.ChatID)
}

// Test_ChatMembersInput_Validate тестирует валидацию входящих параметров
func Test_ChatMembersInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := ChatMembersInput{
			ChatID:    id,
			SubjectID: id,
		}
		return in.Validate()
	})
}

// Test_Members_ChatMembers тестирует получение списка участников чата
func (suite *servicesTestSuite) Test_Members_ChatMembers() {
	suite.Run("чат должен существовать", func() {
		input := ChatMembersInput{
			ChatID:    uuid.NewString(),
			SubjectID: uuid.NewString(),
		}
		members, err := suite.ss.members.ChatMembers(input)
		suite.ErrorIs(err, ErrChatNotExists)
		suite.Empty(members)
	})

	suite.Run("пользователь должен быть участником чата", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать участника в другом чате
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: uuid.NewString(),
		})
		// Запросить список участников чата
		input := ChatMembersInput{
			ChatID:    chat.ID,
			SubjectID: member.UserID,
		}
		members, err := suite.ss.members.ChatMembers(input)
		// Вернется ошибка, потому пользователь не является участником чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
		suite.Empty(members)
	})

	suite.Run("возвращается список участников чата", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать несколько участников в чате
		const membersAllCount = 20
		membersSaved := make([]domain.Member, membersAllCount)
		for i := range membersAllCount {
			// Создать участника в чате
			membersSaved[i] = suite.saveMember(domain.Member{
				ID:     uuid.NewString(),
				UserID: uuid.NewString(),
				ChatID: chat.ID,
			})
		}
		// Запрашивать список будет первый участник
		subjectMember := membersSaved[0]
		// Получить список участников в чате
		input := ChatMembersInput{
			ChatID:    subjectMember.ChatID,
			SubjectID: subjectMember.UserID,
		}
		membersFromRepo, err := suite.ss.members.ChatMembers(input)
		suite.NoError(err)
		if suite.Len(membersFromRepo, membersAllCount) {
			// Сравнить каждого сохраненного участника с ранее созданным
			for i := range membersSaved {
				suite.assertEqualMembers(membersSaved[i], membersFromRepo[i])
			}
		}
	})
}

// Test_LeaveInput_Validate тестирует валидацию входящих параметров
func Test_LeaveInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		in := LeaveChatInput{
			SubjectID: id,
			ChatID:    id,
		}
		return in.Validate()
	})
}

// Test_Members_LeaveChat тестирует выход участника из чата
func (suite *servicesTestSuite) Test_Members_LeaveChat() {
	suite.Run("чат должен существовать", func() {
		// Покинуть чат
		input := LeaveChatInput{
			SubjectID: uuid.NewString(),
			ChatID:    uuid.NewString(),
		}
		err := suite.ss.members.LeaveChat(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
	})

	suite.Run("пользователь должен быть участником чата", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Покинуть чат
		input := LeaveChatInput{
			SubjectID: uuid.NewString(),
			ChatID:    chat.ID,
		}
		err := suite.ss.members.LeaveChat(input)
		// Вернется ошибка, потому что пользователь не участник чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
	})

	suite.Run("пользователь не должен быть главным администратором чата", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Создать участника главного администратора в этом чате
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		})
		// Покинуть чат
		input := LeaveChatInput{
			SubjectID: member.UserID,
			ChatID:    chat.ID,
		}
		err := suite.ss.members.LeaveChat(input)
		// Вернется ошибка, потому что пользователь главный администратор чата
		suite.ErrorIs(err, ErrSubjectUserShouldNotBeChief)
	})

	suite.Run("после выхода пользователь перестает быть участником", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Создать участника в этом чате
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Покинуть чат
		input := LeaveChatInput{
			SubjectID: member.UserID,
			ChatID:    chat.ID,
		}
		err := suite.ss.members.LeaveChat(input)
		suite.Require().NoError(err)
		// Получить список участников чата
		membersFilter := domain.MembersFilter{ID: member.ID}
		members, err := suite.ss.members.MembersRepo.List(membersFilter)
		suite.Require().NoError(err)
		// В чате не осталось участников
		suite.Empty(members)
	})
}

// Test_DeleteMemberInput_Validate тестирует валидацию входящих параметров
func Test_DeleteMemberInput_Validate(t *testing.T) {
	helpers_tests.RunValidateRequiredIDTest(t, func(id string) error {
		input := DeleteMemberInput{
			SubjectID: id,
			ChatID:    id,
			UserID:    id,
		}
		return input.Validate()
	})
}

// Test_Members_DeleteMember тестирует удаление участника чата
func (suite *servicesTestSuite) Test_Members_DeleteMember() {
	suite.Run("нельзя удалить самого себя", func() {
		// Удалить участника
		userID := uuid.NewString()
		input := DeleteMemberInput{
			SubjectID: userID,
			ChatID:    uuid.NewString(),
			UserID:    userID,
		}
		err := suite.ss.members.DeleteMember(input)
		// Вернется ошибка, потому что пользователь пытается удалить самого себя
		suite.ErrorIs(err, ErrMemberCannotDeleteHimself)
	})

	suite.Run("чат должен существовать", func() {
		// Удалить участника
		input := DeleteMemberInput{
			SubjectID: uuid.NewString(),
			ChatID:    uuid.NewString(),
			UserID:    uuid.NewString(),
		}
		err := suite.ss.members.DeleteMember(input)
		// Вернется ошибка, потому что чата не существует
		suite.ErrorIs(err, ErrChatNotExists)
	})

	suite.Run("subject должен быть участником чата", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID: uuid.NewString(),
		})
		// Удалить участника
		input := DeleteMemberInput{
			SubjectID: uuid.NewString(),
			ChatID:    chat.ID,
			UserID:    uuid.NewString(),
		}
		err := suite.ss.members.DeleteMember(input)
		// Вернется ошибка, потому что пользователь не участник чата
		suite.ErrorIs(err, ErrSubjectIsNotMember)
	})

	suite.Run("subject должен быть главным администратором чата", func() {
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
		// Удалить участника
		input := DeleteMemberInput{
			SubjectID: subjectMember.UserID,
			ChatID:    chat.ID,
			UserID:    uuid.NewString(),
		}
		err := suite.ss.members.DeleteMember(input)
		// Вернется ошибка, потому что участник не главный администратор
		suite.ErrorIs(err, ErrSubjectUserIsNotChief)
	})

	suite.Run("user должен быть участником чата", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Создать участника
		member := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		})
		// Удалить участника
		input := DeleteMemberInput{
			SubjectID: member.UserID,
			ChatID:    chat.ID,
			UserID:    uuid.NewString(),
		}
		err := suite.ss.members.DeleteMember(input)
		// Вернется ошибка, потому что удаляемый пользователь не является участником
		suite.ErrorIs(err, ErrUserIsNotMember)
	})

	suite.Run("после удаления участник перестает быть участником", func() {
		// Создать чат
		chat := suite.saveChat(domain.Chat{
			ID:          uuid.NewString(),
			ChiefUserID: uuid.NewString(),
		})
		// Создать участника для удаления
		memberForDelete := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: uuid.NewString(),
			ChatID: chat.ID,
		})
		// Создать участника
		subjectMember := suite.saveMember(domain.Member{
			ID:     uuid.NewString(),
			UserID: chat.ChiefUserID,
			ChatID: chat.ID,
		})
		// Удалить участника
		input := DeleteMemberInput{
			SubjectID: subjectMember.UserID,
			ChatID:    chat.ID,
			UserID:    memberForDelete.UserID,
		}
		err := suite.ss.members.DeleteMember(input)
		suite.Require().NoError(err)
		// Найти удаленного участника
		membersFilter := domain.MembersFilter{ID: memberForDelete.ID}
		members, err := suite.ss.members.MembersRepo.List(membersFilter)
		suite.Require().NoError(err)
		// Такого участника больше нет
		suite.Empty(members)
	})
}
