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

	// проверить является ли пользователь участником чата
	_, err = subjectUserMember(i.MembersRepo, in.SubjectUserID, in.ChatID)
	if err != nil {
		return nil, err
	}

	// Проверить что пользователь является администратором чата
	if chat.ChiefUserID == in.SubjectUserID {
		// Получить все приглашения в этот чат
		invitations, err := getInitationsInThisChat(i.InvitationsRepo, in.ChatID)

		return invitations, err
	} else {
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

	// получить список полученных пользователем приглашений
	invs, err := getInitationsSpecificUser(i.InvitationsRepo, in.UserID)

	return invs, err
}

type SendChatInvitationInput struct {
	ChatID        string
	SubjectUserID string
	UserID        string
}

func (in SendChatInvitationInput) Validate() error {
	if err := uuid.Validate(in.ChatID); err != nil {
		return ErrInvalidChatID
	}
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return ErrInvalidSubjectUserID
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return ErrInvalidUserID
	}

	return nil
}

// SendChatInvitation отправляет приглашения пользователям UserID
func (i *Invitations) SendChatInvitation(in SendChatInvitationInput) error {
	if err := in.Validate(); err != nil {
		return err
	}

	// проверить существование чата
	if _, err := getChat(i.ChatsRepo, in.ChatID); err != nil {
		return err
	}

	// проверить, состоит ли SubjectUserID в чате
	if _, err := subjectUserMember(i.MembersRepo, in.SubjectUserID, in.ChatID); err != nil {
		return err
	}

	// проверить, не состоит ли UserID в чате
	if _, err := userMember(i.MembersRepo, in.UserID, in.ChatID); err == nil {
		return ErrUserAlreadyInChat
	}

	// проверить, существует ли UserID
	if _, err := getUser(i.UsersRepo, in.UserID); err != nil {
		return err
	}

	// проверить, не существет ли приглашение для этого пользователя в этот чат
	if _, err := getInitation(i.InvitationsRepo, in.UserID, in.ChatID); err == nil {
		return ErrUserAlreadyInviteInChat
	}

	// отправить приглашение
	invitation := domain.Invitation{
		ID:            uuid.NewString(),
		SubjectUserID: in.SubjectUserID,
		UserID:        in.UserID,
		ChatID:        in.ChatID,
	}
	err := i.InvitationsRepo.Save(invitation)
	if err != nil {
		return err
	}

	return nil
}

type AcceptInvitationInput struct {
	SubjectUserID string
	ChatID        string
}

func (in AcceptInvitationInput) Validate() error {
	if err := uuid.Validate(in.ChatID); err != nil {
		return ErrInvalidChatID
	}
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return ErrInvalidSubjectUserID
	}

	return nil
}

// AcceptInvitation пинимает приглашения в чат
func (i *Invitations) AcceptInvitation(in AcceptInvitationInput) error {
	// Валидировать входные данные
	if err := in.Validate(); err != nil {
		return err
	}

	// проверить существование приглашения
	invitation, err := getInitation(i.InvitationsRepo, in.SubjectUserID, in.ChatID)
	if err != nil {
		return err
	}

	//проверить существование пользователя
	if _, err := getUser(i.UsersRepo, in.SubjectUserID); err != nil {
		return err
	}

	// проверить, существование чата
	if _, err := getChat(i.ChatsRepo, in.ChatID); err != nil {
		return err
	}

	// создаем участника чата
	member := domain.Member{
		ID:     uuid.NewString(),
		UserID: in.SubjectUserID,
		ChatID: in.ChatID,
	}
	err = i.MembersRepo.Save(member)
	if err != nil {
		return err
	}

	// удаляем приглашение
	err = i.InvitationsRepo.Delete(invitation.ID)
	if err != nil {
		return err
	}

	return nil
}
