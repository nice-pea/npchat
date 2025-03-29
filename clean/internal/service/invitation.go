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
	UserID string
	ChatID string
}

var (
	ErrChatInvitationsInputUserIDValidate = errors.New("некорректный UserID")
	ErrChatInvitationsInputChatIDValidate = errors.New("некорректный ChatID")
	ErrChatInvitationsNoChat              = errors.New("не существует чата с данным ChatID")
	ErrChatInvitationsUserIsNotChief      = errors.New("доступно только для chief этого чата")
)

// Validate валидирует параметры для запроса приглашений конкретного чата
func (in ChatInvitationsInput) Validate() error {
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrChatInvitationsInputChatIDValidate)
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return errors.Join(err, ErrChatInvitationsInputUserIDValidate)
	}
	return nil
}

// ChatInvitations возвращает список приглашений в конкретный чат
func (i *Invitations) ChatInvitations(in ChatInvitationsInput) ([]domain.Invitation, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	chats, err := i.ChatsRepo.List(domain.ChatsFilter{
		IDs: []string{in.ChatID},
	})
	if err != nil {
		return nil, err
	}
	if len(chats) != 1 {
		return nil, ErrChatInvitationsNoChat
	}
	// только 1 чат существует по такому ChatID
	chat := chats[0]
	if chat.ChiefUserID != in.UserID {
		return nil, ErrChatInvitationsUserIsNotChief
	}

	invitations, err := i.InvitationsRepo.List(domain.InvitationsFilter{
		ChatID: chat.ID,
	})

	return invitations, err
}

type UserInvitationsInput struct {
	SubjectUserID string
	UserID        string
}

var (
	ErrUserInvitationsInputSubjectUserIDValidate = errors.New("некорректный SubjectUserID")
	ErrUserInvitationsInputUserIDValidate        = errors.New("некорректный UserID")
	ErrUserInvitationsInputEqualUserIDsValidate  = errors.New("UserID и SubjectUserID не совпадают")
)

func (in UserInvitationsInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrUserInvitationsInputSubjectUserIDValidate)
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return errors.Join(err, ErrUserInvitationsInputUserIDValidate)
	}
	if in.UserID != in.SubjectUserID {
		return ErrUserInvitationsInputEqualUserIDsValidate
	}

	return nil
}

// UserInvitations возвращает список приглашений конкретного пользователя в чаты
func (i *Invitations) UserInvitations(in UserInvitationsInput) ([]domain.Invitation, error) {
	if err := in.Validate(); err != nil {
		return nil, err
	}

	invs, err := i.InvitationsRepo.List(domain.InvitationsFilter{
		ID: in.UserID,
	})

	return invs, err
}
