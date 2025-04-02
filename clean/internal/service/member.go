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

var (
	ErrInvalidSubjectUserID        = errors.New("некорректный SubjectUserID")
	ErrInvalidChatID               = errors.New("некорректный ChatID")
	ErrInvalidUserID               = errors.New("некорректный UserID")
	ErrUserIsNotMember             = errors.New("user не является участником чата")
	ErrSubjectUserIsNotMember      = errors.New("subject user не является участником чата")
	ErrChatNotExists               = errors.New("чата с таким ID не существует")
	ErrMemberCannotDeleteHimself   = errors.New("участник не может удалить самого себя")
	ErrSubjectUserShouldNotBeChief = errors.New("пользователь является главным администратором чата")
	ErrSubjectUserIsNotChief       = errors.New("пользователь не является главным администратором чата")
)

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
	var err error
	if err = in.Validate(); err != nil {
		return nil, err
	}

	// Проверить существование чата
	if _, err = m.chatMustExists(in.ChatID); err != nil {
		return nil, err
	}

	// Получить список участников
	var members []domain.Member
	if members, err = m.chatMembers(in.ChatID); err != nil {
		return nil, err
	}

	// Проверить что пользователь является участником чата
	if len(members) == 0 {
		return nil, ErrSubjectUserIsNotMember
	}
	for i, member := range members {
		if member.UserID == in.SubjectUserID {
			break
		} else if i == len(members)-1 {
			return nil, ErrSubjectUserIsNotMember
		}
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
	var err error
	if err = in.Validate(); err != nil {
		return err
	}

	// Проверить существование чата
	var chat domain.Chat
	if chat, err = m.chatMustExists(in.ChatID); err != nil {
		return err
	}

	// Пользователь должен быть участником чата
	var subjectMember domain.Member
	if subjectMember, err = m.subjectUserMustBeMember(in.SubjectUserID, chat.ID); err != nil {
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
	var err error
	if err = in.Validate(); err != nil {
		return err
	}

	// Проверить попытку удалить самого себя
	if in.UserID == in.SubjectUserID {
		return ErrMemberCannotDeleteHimself
	}

	// Проверить существование чата
	var chat domain.Chat
	if chat, err = m.chatMustExists(in.ChatID); err != nil {
		return err
	}

	// Пользователь должен быть участником чата
	if _, err = m.subjectUserMustBeMember(in.SubjectUserID, chat.ID); err != nil {
		return err
	}

	// Пользователь должен быть главным администратором
	if chat.ChiefUserID != in.SubjectUserID {
		return ErrSubjectUserIsNotChief
	}

	// Удаляемый участник должен существовать
	var member domain.Member
	if member, err = m.userMustBeMember(in.UserID, chat.ID); err != nil {
		return err
	}

	// Удалить участника
	if err = m.MembersRepo.Delete(member.ID); err != nil {
		return err
	}

	return nil
}

// Проверить существование чата
func (m *Members) chatMustExists(chatID string) (domain.Chat, error) {
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

// Пользователь должен быть участником чата
func (m *Members) userMustBeMember(userID, chatID string) (domain.Member, error) {
	membersFilter := domain.MembersFilter{
		UserID: userID,
		ChatID: chatID,
	}
	members, err := m.MembersRepo.List(membersFilter)
	if err != nil {
		return domain.Member{}, err
	}
	if len(members) != 1 {
		return domain.Member{}, ErrUserIsNotMember
	}

	return members[0], nil
}

// Пользователь должен быть участником чата
func (m *Members) subjectUserMustBeMember(subjectUserID, chatID string) (domain.Member, error) {
	membersFilter := domain.MembersFilter{
		UserID: subjectUserID,
		ChatID: chatID,
	}
	members, err := m.MembersRepo.List(membersFilter)
	if err != nil {
		return domain.Member{}, err
	}
	if len(members) != 1 {
		return domain.Member{}, ErrSubjectUserIsNotMember
	}

	return members[0], nil
}
