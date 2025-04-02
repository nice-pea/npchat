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
	chatsFilter := domain.ChatsFilter{
		IDs: []string{in.ChatID},
	}
	chats, err := m.ChatsRepo.List(chatsFilter)
	if err != nil {
		return nil, err
	}
	if len(chats) != 1 {
		return nil, ErrChatNotExists
	}

	// Получить список участников
	membersFilter := domain.MembersFilter{
		ChatID: in.ChatID,
	}
	members, err := m.MembersRepo.List(membersFilter)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, ErrSubjectUserIsNotMember
	}

	// Проверить что пользователь является участником чата
	for i, member := range members {
		if member.UserID == in.SubjectUserID {
			break
		} else if i == len(members)-1 {
			return nil, ErrSubjectUserIsNotMember
		}
	}

	return members, nil
}

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
	chatsFilter := domain.ChatsFilter{
		IDs: []string{in.ChatID},
	}
	chats, err := m.ChatsRepo.List(chatsFilter)
	if err != nil {
		return err
	}
	if len(chats) != 1 {
		return ErrChatNotExists
	}

	// Пользователь должен быть участником чата
	membersFilter := domain.MembersFilter{
		UserID: in.SubjectUserID,
		ChatID: in.ChatID,
	}
	members, err := m.MembersRepo.List(membersFilter)
	if err != nil {
		return err
	}
	if len(members) != 1 {
		return ErrSubjectUserIsNotMember
	}

	// Пользователь не должен быть главным администратором
	if members[0].UserID == chats[0].ChiefUserID {
		return ErrSubjectUserShouldNotBeChief
	}

	// Удалить пользователя из чата
	if err = m.MembersRepo.Delete(members[0].ID); err != nil {
		return err
	}

	return nil
}

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
	chatsFilter := domain.ChatsFilter{
		IDs: []string{in.ChatID},
	}
	chats, err := m.ChatsRepo.List(chatsFilter)
	if err != nil {
		return err
	}
	if len(chats) != 1 {
		return ErrChatNotExists
	}

	// Пользователь должен быть участником чата
	membersFilter := domain.MembersFilter{
		UserID: in.SubjectUserID,
		ChatID: in.ChatID,
	}
	members, err := m.MembersRepo.List(membersFilter)
	if err != nil {
		return err
	}
	if len(members) != 1 {
		return ErrSubjectUserIsNotMember
	}

	// Пользователь должен быть главным администратором
	if chats[0].ChiefUserID != in.SubjectUserID {
		return ErrSubjectUserIsNotChief
	}

	// Удаляемый участник должен существовать
	membersFilter = domain.MembersFilter{
		UserID: in.UserID,
		ChatID: in.ChatID,
	}
	members, err = m.MembersRepo.List(membersFilter)
	if err != nil {
		return err
	}
	if len(members) != 1 {
		return ErrUserIsNotMember
	}

	// Удалить участника
	if err = m.MembersRepo.Delete(members[0].ID); err != nil {
		return err
	}

	return nil
}
