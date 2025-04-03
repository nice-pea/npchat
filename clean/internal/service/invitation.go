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
	UsersRepo       domain.UsersRepository
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
	chat, err := getChat(i.ChatsRepo, in.ChatID)
	if err != nil {
		return nil, err
	}

	// Проверить что пользователь является администратором чата
	if chat.ChiefUserID == in.SubjectUserID {
		// Получить все приглашения в этот чат
		invitations, err := getInitationsInThisChat(i.InvitationsRepo, in.ChatID)

		return invitations, err
	} else {
		// проверить является ли пользователь участником чата
		_, err := subjectUserMember(i.MembersRepo, in.SubjectUserID, in.ChatID)
		if err != nil {
			return nil, err
		}

		// получить список приглашений конкретного пользователя
		invitations, err := getInitationsSpecificUserInThisChat(i.InvitationsRepo, in.SubjectUserID, in.ChatID)

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
	_, err := getUser(i.UsersRepo, in.UserID)
	if err != nil {
		return nil, err
	}

	// получить список приглашений
	invs, err := getInitationsSpecificUser(i.InvitationsRepo, in.UserID)

	return invs, err
}

// getInitationsInThisChat возвращает список всех приглашений в конкретный чат
func getInitationsInThisChat(invitationsRepo domain.InvitationsRepository, chatId string) ([]domain.Invitation, error) {
	invitations, err := invitationsRepo.List(domain.InvitationsFilter{
		ChatID: chatId,
	})
	return invitations, err
}

// getInitationsSpecificUserInThisChat возвращает список приглашений конкретного пользователя в конкретный чат
func getInitationsSpecificUserInThisChat(invitationsRepo domain.InvitationsRepository, subjectUserId, chatId string) ([]domain.Invitation, error) {
	invitations, err := invitationsRepo.List(domain.InvitationsFilter{
		SubjectUserID: subjectUserId,
		ChatID:        chatId,
	})
	return invitations, err
}

// getInitationsSpecificUser возвращает список приглашений конкретного пользователя в чаты
func getInitationsSpecificUser(invitationsRepo domain.InvitationsRepository, userId string) ([]domain.Invitation, error) {
	invitations, err := invitationsRepo.List(domain.InvitationsFilter{
		UserID: userId,
	})
	return invitations, err
}
