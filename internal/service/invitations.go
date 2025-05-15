package service

import (
	"errors"
	"slices"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Invitations сервис, объединяющий случаи использования(юзкейсы) в контексте сущности
type Invitations struct {
	ChatsRepo       domain.ChatsRepository
	MembersRepo     domain.MembersRepository
	InvitationsRepo domain.InvitationsRepository
	UsersRepo       domain.UsersRepository
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

// ChatInvitations возвращает список приглашений в конкретный чат.
// Если SubjectUserID является администратором, то возвращается все приглашения в данный чат,
// иначе только те приглашения, которые отправил именно пользователь.
func (i *Invitations) ChatInvitations(in ChatInvitationsInput) ([]domain.Invitation, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return nil, err
	}

	// Проверить существование чата
	chat, err := getChat(i.ChatsRepo, in.ChatID)
	if err != nil {
		return nil, err
	}

	// Проверить является ли пользователь участником чата
	_, err = subjectUserMember(i.MembersRepo, in.SubjectUserID, in.ChatID)
	if err != nil {
		return nil, err
	}

	// Проверить, что пользователь является администратором чата
	if chat.ChiefUserID == in.SubjectUserID {
		// Получить все приглашения в этот чат
		invitations, err := i.InvitationsRepo.List(domain.InvitationsFilter{
			ChatID: in.ChatID,
		})

		return invitations, err
	} else {
		// Получить список приглашений конкретного пользователя
		invitations, err := i.InvitationsRepo.List(domain.InvitationsFilter{
			SubjectUserID: in.SubjectUserID,
			ChatID:        in.ChatID,
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
		return nil, ErrUnauthorizedInvitationsView
	}

	// Пользователь должен существовать
	_, err := getUser(i.UsersRepo, in.UserID)
	if err != nil {
		return nil, err
	}

	// Получить список полученных пользователем приглашений
	invitations, err := i.InvitationsRepo.List(domain.InvitationsFilter{
		UserID: in.UserID,
	})

	return invitations, err
}

type SendInvitationInput struct {
	SubjectUserID string
	ChatID        string
	UserID        string
}

func (in SendInvitationInput) Validate() error {
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

// SendInvitation отправляет приглашения пользователю от участника чата
func (i *Invitations) SendInvitation(in SendInvitationInput) (domain.Invitation, error) {
	if err := in.Validate(); err != nil {
		return domain.Invitation{}, err
	}

	// Проверить существование чата
	if _, err := getChat(i.ChatsRepo, in.ChatID); err != nil {
		return domain.Invitation{}, err
	}

	// Проверить, состоит ли SubjectUserID в чате
	if _, err := subjectUserMember(i.MembersRepo, in.SubjectUserID, in.ChatID); err != nil {
		return domain.Invitation{}, err
	}

	// Проверить, не состоит ли UserID в чате
	if _, err := userMember(i.MembersRepo, in.UserID, in.ChatID); err == nil {
		return domain.Invitation{}, ErrUserIsAlreadyInChat
	}

	// Проверить, существует ли UserID
	if _, err := getUser(i.UsersRepo, in.UserID); err != nil {
		return domain.Invitation{}, err
	}

	// Проверить, не существует ли приглашение для этого пользователя в этот чат
	if err := invitationMustNotExist(i.InvitationsRepo, in.UserID, in.ChatID); err != nil {
		return domain.Invitation{}, ErrUserIsAlreadyInvited
	}

	// Отправить приглашение
	invitation := domain.Invitation{
		ID:            uuid.NewString(),
		SubjectUserID: in.SubjectUserID,
		UserID:        in.UserID,
		ChatID:        in.ChatID,
	}
	err := i.InvitationsRepo.Save(invitation)
	if err != nil {
		return domain.Invitation{}, err
	}

	return invitation, nil
}

type AcceptInvitationInput struct {
	SubjectUserID string
	InvitationID  string
}

func (in AcceptInvitationInput) Validate() error {
	if err := uuid.Validate(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return ErrInvalidSubjectUserID
	}

	return nil
}

// AcceptInvitation добавляет пользователя в чат, путем принятия приглашения
func (i *Invitations) AcceptInvitation(in AcceptInvitationInput) error {
	var err error

	// Валидировать входные данные
	if err = in.Validate(); err != nil {
		return err
	}

	// Проверить существование приглашения
	invitation, err := getInvitationByID(i.InvitationsRepo, in.InvitationID)
	if err != nil {
		return err
	}

	// Проверить существование пользователя
	user, err := getUser(i.UsersRepo, in.SubjectUserID)
	if err != nil {
		return err
	}

	// Приглашение должно быть направлено пользователю
	if invitation.UserID != user.ID {
		return ErrSubjectUserNotAllowed
	}

	// Создаем участника чата
	member := domain.Member{
		ID:     uuid.NewString(),
		UserID: in.SubjectUserID,
		ChatID: invitation.ChatID,
	}
	err = i.MembersRepo.Save(member)
	if err != nil {
		return err
	}

	// Удаляем приглашение
	err = i.InvitationsRepo.Delete(invitation.ID)
	if err != nil {
		return err
	}

	return nil
}

type CancelInvitationInput struct {
	SubjectUserID string
	InvitationID  string
}

func (in CancelInvitationInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return ErrInvalidSubjectUserID
	}
	if err := uuid.Validate(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}

	return nil
}

// CancelInvitation отменяет приглашение
func (i *Invitations) CancelInvitation(in CancelInvitationInput) error {
	// Валидировать входные данные
	if err := in.Validate(); err != nil {
		return err
	}

	// Проверить существование приглашения
	invitation, err := getInvitationByID(i.InvitationsRepo, in.InvitationID)
	if err != nil {
		return err
	}

	// Найти чат
	chat, err := getChat(i.ChatsRepo, invitation.ChatID)
	if err != nil && !errors.Is(err, ErrChatNotExists) {
		return err
	}

	if in.SubjectUserID == invitation.SubjectUserID {
		// Проверить, существование участника чата
		_, err := subjectUserMember(i.MembersRepo, invitation.SubjectUserID, chat.ID)
		if err != nil {
			return err
		}
	}

	// Список тех кто может отменять приглашение
	allowedSubjects := []string{
		chat.ChiefUserID,         // Главный администратор
		invitation.SubjectUserID, // Пригласивший
		invitation.UserID,        // Приглашаемый
	}
	// Проверить, может ли пользователь отменить приглашение
	if !slices.Contains(allowedSubjects, in.SubjectUserID) {
		return ErrSubjectUserNotAllowed
	}

	err = i.InvitationsRepo.Delete(invitation.ID)
	if err != nil {
		return err
	}

	return nil
}

// getInvitation возвращает приглашение в конкретный чат
func getInvitation(invitationsRepo domain.InvitationsRepository, userId, chatId string) (domain.Invitation, error) {
	invitations, err := invitationsRepo.List(domain.InvitationsFilter{
		UserID: userId,
		ChatID: chatId,
	})
	if err != nil {
		return domain.Invitation{}, err
	}
	if len(invitations) != 1 {
		return domain.Invitation{}, ErrInvitationNotExists
	}

	return invitations[0], nil
}

// getInvitation возвращает приглашение в конкретный чат
func getInvitationByID(invitationsRepo domain.InvitationsRepository, id string) (domain.Invitation, error) {
	invitations, err := invitationsRepo.List(domain.InvitationsFilter{
		ID: id,
	})
	if err != nil {
		return domain.Invitation{}, err
	}
	if len(invitations) != 1 {
		return domain.Invitation{}, ErrInvitationNotExists
	}

	return invitations[0], nil
}

// invitationMustNotExist возвращает ошибку ErrUserIsAlreadyInvited, если приглашение по таким фильтрам существует
func invitationMustNotExist(invitationsRepo domain.InvitationsRepository, userId, chatId string) error {
	invitations, err := invitationsRepo.List(domain.InvitationsFilter{
		UserID: userId,
		ChatID: chatId,
	})
	if err != nil {
		return err
	}
	if len(invitations) > 0 {
		return ErrUserIsAlreadyInvited
	}

	return nil
}
