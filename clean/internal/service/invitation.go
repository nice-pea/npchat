package service

import (
	"errors"

	"github.com/google/uuid"
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Invitations сервис объединяющий случаи использования(юзкейсы) в контексте сущности
type Invitations struct {
	ChatsRepo       domain.ChatsRepository
	MembersRepo     domain.MembersRepository
	InvitationsRepo domain.InvitationsRepository
	History         History
}

// ChatInvitationsInput параметры для запроса приглашений конкретного чата
type ChatInvitationsInput struct {
	SubjectUserID string
	ChatID        string
}

// Validate валидирует параметры для запроса приглашений конкретного чата
func (in ChatInvitationsInput) Validate() error {
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvalidSubjectUserID)
	}
	return nil
}

// ChatInvitations возвращает список приглашений в конкретный чат
// если SubjectUserID является администратором то возвращается все приглашения в данный чат
// иначе только те приглашения которые отправил именно пользователь
func (i *Invitations) ChatInvitations(in ChatInvitationsInput) ([]domain.Invitation, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// проверить существование чата
	chats, err := i.ChatsRepo.List(domain.ChatsFilter{
		IDs: []string{in.ChatID},
	})
	if err != nil {
		return nil, err
	}
	if len(chats) != 1 {
		return nil, ErrChatNotExists
	}
	chat := chats[0]

	// Проверить что пользователь является администратором чата
	if chat.ChiefUserID == in.SubjectUserID {
		// Получить все приглашения в этот чат
		invitations, err := i.InvitationsRepo.List(domain.InvitationsFilter{
			ChatID: chat.ID,
		})

		return invitations, err
	} else {
		// проверить является ли пользователь участником чата
		members, err := i.MembersRepo.List(domain.MembersFilter{
			UserID: in.SubjectUserID,
			ChatID: in.ChatID,
		})
		if err != nil {
			return nil, err
		}
		if len(members) != 1 {
			return nil, ErrUserIsNotMember
		}

		// получить список приглашений конкретного пользователя
		invitations, err := i.InvitationsRepo.List(domain.InvitationsFilter{
			ChatID:        chat.ID,
			SubjectUserID: in.SubjectUserID,
		})

		return invitations, err
	}

}

type UserInvitationsInput struct {
	SubjectUserID string
	UserID        string
}

func (in UserInvitationsInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvalidSubjectUserID)
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}

	return nil
}

// UserInvitations возвращает список приглашений конкретного пользователя в чаты
func (i *Invitations) UserInvitations(in UserInvitationsInput) ([]domain.Invitation, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Пользователь должен видеть только свои приглашения
	if in.UserID != in.SubjectUserID {
		return nil, ErrCannotViewSomeoneElseChats
	}

	// Пользователь должен существовать

	// получить список приглашений
	invs, err := i.InvitationsRepo.List(domain.InvitationsFilter{
		UserID: in.UserID,
	})

	return invs, err
}

// chat возвращает чат либо ошибку ErrChatNotExists
func (i *Invitations) chat(chatID string) (domain.Chat, error) {
	chatsFilter := domain.ChatsFilter{
		IDs: []string{chatID},
	}
	chats, err := i.ChatsRepo.List(chatsFilter)
	if err != nil {
		return domain.Chat{}, err
	}
	if len(chats) != 1 {
		return domain.Chat{}, ErrChatNotExists
	}

	return chats[0], nil
}

// Получить список участников
func (i *Invitations) chatMembers(chatID string) ([]domain.Member, error) {
	membersFilter := domain.MembersFilter{
		ChatID: chatID,
	}
	members, err := i.MembersRepo.List(membersFilter)
	if err != nil {
		return nil, err
	}

	return members, nil
}

// UserMember вернет участника либо ошибку ErrUserIsNotMember
func (i *Invitations) userMember(userID, chatID string) (domain.Member, error) {
	return i.memberOrErr(userID, chatID, ErrUserIsNotMember)
}

// subjectUserMember вернет участника либо ошибку ErrSubjectUserIsNotMember
func (i *Invitations) subjectUserMember(subjectUserID, chatID string) (domain.Member, error) {
	return i.memberOrErr(subjectUserID, chatID, ErrSubjectUserIsNotMember)
}

// memberOrErr возвращает участника чата по userID, chatID.
// Вернет errOnNotExists ошибку если участника не будет существовать.
func (i *Invitations) memberOrErr(userID, chatID string, errOnNotExists error) (domain.Member, error) {
	membersFilter := domain.MembersFilter{
		UserID: userID,
		ChatID: chatID,
	}
	members, err := i.MembersRepo.List(membersFilter)
	if err != nil {
		return domain.Member{}, err
	}
	if len(members) != 1 {
		return domain.Member{}, errOnNotExists
	}

	return members[0], nil
}
