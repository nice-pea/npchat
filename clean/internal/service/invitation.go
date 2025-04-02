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

var (
	ErrChatInvitationsInputSubjectUserIDValidate = errors.New("некорректный SubjectUserID")
	ErrChatInvitationsInputChatIDValidate        = errors.New("некорректный ChatID")
	ErrChatInvitationsNoChat                     = errors.New("не существует чата с данным ChatID")
	ErrChatInvitationsUserIsNotChief             = errors.New("доступно только для администратора этого чата")
	ErrChatInvitationsUserIsNotMember            = errors.New("пользователь не является участником чата")
)

// Validate валидирует параметры для запроса приглашений конкретного чата
func (in ChatInvitationsInput) Validate() error {
	if err := uuid.Validate(in.ChatID); err != nil {
		return errors.Join(err, ErrChatInvitationsInputChatIDValidate)
	}
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrChatInvitationsInputSubjectUserIDValidate)
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
		return nil, ErrChatInvitationsNoChat
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
			return nil, ErrChatInvitationsUserIsNotMember
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

var (
	ErrUserInvitationsInputSubjectUserIDValidate = errors.New("некорректный SubjectUserID")
	ErrUserInvitationsInputUserIDValidate        = errors.New("некорректный UserID")
	ErrUserInvitationsInputEqualUserIDsValidate  = errors.New("доступно только самому пользователю")
)

func (in UserInvitationsInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrUserInvitationsInputSubjectUserIDValidate)
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return errors.Join(err, ErrUserInvitationsInputUserIDValidate)
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
		return nil, ErrUserInvitationsInputEqualUserIDsValidate
	}

	// Пользователь должен существовать
	

	// получить список приглашений
	invs, err := i.InvitationsRepo.List(domain.InvitationsFilter{
		UserID: in.UserID,
	})

	return invs, err
}
