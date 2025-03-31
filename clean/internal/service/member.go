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

type ChatMembersInput struct {
	SubjectUserID string
	ChatID        string
}

var (
	ErrChatMembersInputSubjectUserIDValidate = errors.New("некорректный SubjectUserID")
	ErrChatMembersInputChatIDValidate        = errors.New("некорректный ChatID")
)

// Validate валидирует значение отдельно каждого параметры
func (in ChatMembersInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrChatMembersInputSubjectUserIDValidate)
	}
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrChatMembersInputChatIDValidate)
	}

	return nil
}

var (
	ErrChatMembersChatNotExists   = errors.New("чата с таким ID не существует")
	ErrChatMembersUserIsNotMember = errors.New("пользователь не является участником чата")
)

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
		return nil, ErrChatMembersChatNotExists
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
		return nil, ErrChatMembersUserIsNotMember
	}

	// Проверить что пользователь является участником чата
	for i, member := range members {
		if member.UserID == in.SubjectUserID {
			break
		} else if i == len(members)-1 {
			return nil, ErrChatMembersUserIsNotMember
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
		return errors.Join(err, ErrChatMembersInputSubjectUserIDValidate)
	}
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrChatMembersInputChatIDValidate)
	}

	return nil
}

var (
	ErrMembersLeaveChatNotExists            = errors.New("чата с таким ID не существует")
	ErrMembersLeaveChatUserIsNotMember      = errors.New("пользователь не является участником чата")
	ErrMembersLeaveChatUserShouldNotBeChief = errors.New("пользователь является главным администратором чата")
)

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
		return ErrMembersLeaveChatNotExists
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
		return ErrMembersLeaveChatUserIsNotMember
	}

	// Пользователь не должен быть главным администратором
	if members[0].UserID == chats[0].ChiefUserID {
		return ErrMembersLeaveChatUserShouldNotBeChief
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

var (
	ErrDeleteMemberInputSubjectUserIDValidate = errors.New("некорректный SubjectUserID")
	ErrDeleteMemberInputChatIDValidate        = errors.New("некорректный ChatID")
	ErrDeleteMemberInputUserIDValidate        = errors.New("некорректный UserID")
)

// Validate валидирует значение отдельно каждого параметры
func (in DeleteMemberInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrDeleteMemberInputSubjectUserIDValidate)
	}
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrDeleteMemberInputChatIDValidate)
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return errors.Join(err, ErrDeleteMemberInputUserIDValidate)
	}

	return nil
}

var (
	ErrMembersDeleteMemberCannotDeleteHimself    = errors.New("пользователь не может удалить самого себя")
	ErrMembersDeleteMemberMemberIsNotExists      = errors.New("участника не существует чата")
	ErrMembersDeleteMemberSubjectUserIsNotMember = errors.New("пользователь не является участником чата")
	ErrMembersDeleteMemberSubjectUserIsNotChief  = errors.New("пользователь не является главным администратором чата")
)

// DeleteMember удаляет участника чата
func (m *Members) DeleteMember(in DeleteMemberInput) error {
	// Валидировать параметры
	var err error
	if err = in.Validate(); err != nil {
		return err
	}

	// Проверить попытку удалить самого себя
	if in.UserID == in.SubjectUserID {
		return ErrMembersDeleteMemberCannotDeleteHimself
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
		return ErrMembersLeaveChatNotExists
	}

	// Пользователь должен быть главным администратором
	if chats[0].ChiefUserID != in.SubjectUserID {
		return ErrMembersDeleteMemberSubjectUserIsNotChief
	}

	// Удаляемый участник должен существовать
	membersFilter := domain.MembersFilter{
		UserID: in.UserID,
		ChatID: in.ChatID,
	}
	members, err := m.MembersRepo.List(membersFilter)
	if err != nil {
		return err
	}
	if len(members) != 1 {
		return ErrMembersDeleteMemberMemberIsNotExists
	}
	memberForDelete := members[0]

	// Пользователь должен быть участником чата
	membersFilter = domain.MembersFilter{
		UserID: in.SubjectUserID,
		ChatID: in.ChatID,
	}
	if members, err = m.MembersRepo.List(membersFilter); err != nil {
		return err
	}
	if len(members) != 1 {
		return ErrMembersDeleteMemberSubjectUserIsNotMember
	}

	return m.MembersRepo.Delete(memberForDelete.ID)
}
