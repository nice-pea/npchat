package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type Members struct {
	MembersRepo domain.MembersRepository
	ChatsRepo   domain.ChatsRepository
}



// ChatMembersInput входящие параметры
type ChatMembersInput struct {
	SubjectUserID string
	ChatID        string
}

// Validate валидирует значение отдельно каждого параметры
func (in ChatMembersInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}

	return nil
}

// ChatMembers возвращает список участников чата
func (m *Members) ChatMembers(in ChatMembersInput) ([]domain.Member, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Проверить существование чата
	if _, err := m.chat(in.ChatID); err != nil {
		return nil, err
	}

	// Пользователь должен быть участником чата
	if _, err := m.subjectUserMember(in.SubjectUserID, in.ChatID); err != nil {
		return nil, err
	}

	// Получить список участников
	members, err := m.chatMembers(in.ChatID)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// LeaveChatInput входящие параметры
type LeaveChatInput struct {
	SubjectUserID string
	ChatID        string
}

// Validate валидирует значение отдельно каждого параметры
func (in LeaveChatInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvalidSubjectUserID)
	}
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}

	return nil
}

// LeaveChat удаляет участника из чата
func (m *Members) LeaveChat(in LeaveChatInput) error {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return err
	}

	// Проверить существование чата
	chat, err := m.chat(in.ChatID)
	if err != nil {
		return err
	}

	// Пользователь должен быть участником чата
	subjectMember, err := m.subjectUserMember(in.SubjectUserID, chat.ID)
	if err != nil {
		return err
	}

	// Пользователь не должен быть главным администратором
	if in.SubjectUserID == chat.ChiefUserID {
		return ErrSubjectUserShouldNotBeChief
	}

	// Удалить пользователя из чата
	if err = m.MembersRepo.Delete(subjectMember.ID); err != nil {
		return err
	}

	return nil
}

// DeleteMemberInput входящие параметры
type DeleteMemberInput struct {
	SubjectUserID string
	ChatID        string
	UserID        string
}

// Validate валидирует значение отдельно каждого параметры
func (in DeleteMemberInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvalidSubjectUserID)
	}
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}

	return nil
}

// DeleteMember удаляет участника чата
func (m *Members) DeleteMember(in DeleteMemberInput) error {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return err
	}

	// Проверить попытку удалить самого себя
	if in.UserID == in.SubjectUserID {
		return ErrMemberCannotDeleteHimself
	}

	// Проверить существование чата
	chat, err := m.chat(in.ChatID)
	if err != nil {
		return err
	}

	// Пользователь должен быть участником чата
	if _, err = m.subjectUserMember(in.SubjectUserID, in.ChatID); err != nil {
		return err
	}

	// Пользователь должен быть главным администратором
	if chat.ChiefUserID != in.SubjectUserID {
		return ErrSubjectUserIsNotChief
	}

	// Удаляемый участник должен существовать
	member, err := m.userMember(in.UserID, in.ChatID)
	if err != nil {
		return err
	}

	// Удалить участника
	if err = m.MembersRepo.Delete(member.ID); err != nil {
		return err
	}

	return nil
}

// chat возвращает чат либо ошибку ErrChatNotExists
func (m *Members) chat(chatID string) (domain.Chat, error) {
	chatsFilter := domain.ChatsFilter{
		IDs: []string{chatID},
	}
	chats, err := m.ChatsRepo.List(chatsFilter)
	if err != nil {
		return domain.Chat{}, err
	}
	if len(chats) != 1 {
		return domain.Chat{}, ErrChatNotExists
	}

	return chats[0], nil
}

// Получить список участников
func (m *Members) chatMembers(chatID string) ([]domain.Member, error) {
	membersFilter := domain.MembersFilter{
		ChatID: chatID,
	}
	members, err := m.MembersRepo.List(membersFilter)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// subjectUserMember вернет участника либо ошибку ErrUserIsNotMember
func (m *Members) userMember(userID, chatID string) (domain.Member, error) {
	return m.memberOrErr(userID, chatID, ErrUserIsNotMember)
}

// subjectUserMember вернет участника либо ошибку ErrSubjectUserIsNotMember
func (m *Members) subjectUserMember(subjectUserID, chatID string) (domain.Member, error) {
	return m.memberOrErr(subjectUserID, chatID, ErrSubjectUserIsNotMember)
}

// memberOrErr возвращает участника чата по userID, chatID.
// Вернет errOnNotExists ошибку если участника не будет существовать.
func (m *Members) memberOrErr(userID, chatID string, errOnNotExists error) (domain.Member, error) {
	membersFilter := domain.MembersFilter{
		UserID: userID,
		ChatID: chatID,
	}
	members, err := m.MembersRepo.List(membersFilter)
	if err != nil {
		return domain.Member{}, err
	}
	if len(members) != 1 {
		return domain.Member{}, errOnNotExists
	}

	return members[0], nil
}
