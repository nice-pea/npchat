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
	ChatID        string
	SubjectUserID string
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
		return nil, errors.New("чата с таким ID не существует")
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
		return nil, errors.New("пользователь не является участником чата")
	}

	// Проверить что пользователь является участником чата
	for i, member := range members {
		if member.UserID == in.SubjectUserID {
			break
		} else if i == len(members)-1 {
			return nil, errors.New("пользователь не является участником чата")
		}
	}

	return members, nil
}
