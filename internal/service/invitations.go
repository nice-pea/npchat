package service

import (
	"errors"
	"slices"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
)

// ChatInvitationsInput параметры для запроса приглашений конкретного чата
type ChatInvitationsInput struct {
	SubjectID string
	ChatID    string
}

// Validate валидирует параметры для запроса приглашений конкретного чата
func (in ChatInvitationsInput) Validate() error {
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// ChatInvitationsOutput результат запроса приглашений конкретного чата
type ChatInvitationsOutput struct {
	Invitations []chatt.Invitation
}

// ChatInvitations возвращает список приглашений в конкретный чат.
// Если SubjectID является администратором, то возвращается все приглашения в данный чат,
// иначе только те приглашения, которые отправил именно пользователь.
func (c *Chats) ChatInvitations(in ChatInvitationsInput) (ChatInvitationsOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return ChatInvitationsOutput{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return ChatInvitationsOutput{}, err
	}

	// Проверить является ли пользователь участником чата
	if !chat.HasParticipant(in.SubjectID) {
		return ChatInvitationsOutput{}, ErrSubjectIsNotMember
	}

	// Сохранить сначала все приглашения
	invitations := chat.Invitations

	// Если пользователь не является администратором,
	// то оставить только те приглашения, которые отправил именно пользователь.
	if chat.ChiefID != in.SubjectID {
		invitations = chat.SubjectInvitations(in.SubjectID)
	}

	return ChatInvitationsOutput{
		Invitations: invitations,
	}, err
}

type ReceivedInvitationsInput struct {
	SubjectID string
}

func (in ReceivedInvitationsInput) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// ReceivedInvitationsOutput входящие параметры
type ReceivedInvitationsOutput struct {
	// ChatsInvitations карта приглашений, где ключ - chatID, значение - приглашение
	ChatsInvitations map[string]chatt.Invitation
}

// ReceivedInvitations возвращает список приглашений конкретного пользователя в чаты
func (c *Chats) ReceivedInvitations(in ReceivedInvitationsInput) (ReceivedInvitationsOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return ReceivedInvitationsOutput{}, err
	}

	// Найти чат
	chats, err := c.Repo.List(chatt.Filter{
		InvitationRecipientID: in.SubjectID,
	})
	if err != nil {
		return ReceivedInvitationsOutput{}, err
	}

	// Если нет чатов, вернут пустой список
	if len(chats) == 0 {
		return ReceivedInvitationsOutput{}, nil
	}

	// Собрать приглашения, полученные пользователем
	invitations := make(map[string]chatt.Invitation, len(chats))
	for _, chat := range chats {
		invitations[chat.ID], _ = chat.RecipientInvitation(in.SubjectID)
	}

	return ReceivedInvitationsOutput{
		ChatsInvitations: invitations,
	}, err
}

type SendInvitationInput struct {
	SubjectID string
	ChatID    string
	UserID    string
}

func (in SendInvitationInput) Validate() error {
	if err := domain.ValidateID(in.ChatID); err != nil {
		return ErrInvalidChatID
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}
	if err := domain.ValidateID(in.UserID); err != nil {
		return ErrInvalidUserID
	}

	return nil
}

type SendInvitationOutput struct {
	Invitation chatt.Invitation
}

// SendInvitation отправляет приглашения пользователю от участника чата
func (c *Chats) SendInvitation(in SendInvitationInput) (SendInvitationOutput, error) {
	if err := in.Validate(); err != nil {
		return SendInvitationOutput{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return SendInvitationOutput{}, err
	}

	// Создать приглашение
	inv, err := chatt.NewInvitation(in.SubjectID, in.UserID)
	if err != nil {
		return SendInvitationOutput{}, err
	}

	// Добавить приглашение в чат
	if err = chat.AddInvitation(inv); err != nil {
		return SendInvitationOutput{}, err
	}

	// Сохранить чат в репозиторий
	if err = c.Repo.Upsert(chat); err != nil {
		return SendInvitationOutput{}, err
	}

	return SendInvitationOutput{
		Invitation: inv,
	}, nil
}

type AcceptInvitationInput struct {
	SubjectID    string
	InvitationID string
}

func (in AcceptInvitationInput) Validate() error {
	if err := domain.ValidateID(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}

	return nil
}

// AcceptInvitation добавляет пользователя в чат, путем принятия приглашения
func (c *Chats) AcceptInvitation(in AcceptInvitationInput) error {
	// Валидировать входные данные
	if err := in.Validate(); err != nil {
		return err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{
		InvitationID: in.InvitationID,
	})
	if err != nil {
		return err
	}

	// Удаляем приглашение из чата
	if err := chat.RemoveInvitation(in.InvitationID); err != nil {
		return err
	}

	// Создаем участника чата
	participant, err := chatt.NewParticipant(in.SubjectID)
	if err != nil {
		return err
	}

	// Добавить участника в чат
	if err := chat.AddParticipant(participant); err != nil {
		return err
	}

	// Сохранить чат в репозиторий
	if err := c.Repo.Upsert(chat); err != nil {
		return err
	}

	return nil
}

type CancelInvitationInput struct {
	SubjectID    string
	InvitationID string
}

func (in CancelInvitationInput) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}
	if err := domain.ValidateID(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}

	return nil
}

// CancelInvitation отменяет приглашение
func (c *Chats) CancelInvitation(in CancelInvitationInput) error {
	// Валидировать входные данные
	if err := in.Validate(); err != nil {
		return err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{
		InvitationID: in.InvitationID,
	})
	if err != nil {
		return err
	}

	// Достать приглашение из чата
	invitation, err := chat.Invitation(in.InvitationID)
	if err != nil {
		return err
	}

	if in.SubjectID == invitation.SubjectID {
		// Проверить, существование участника чата
		if !chat.HasParticipant(invitation.SubjectID) {
			return ErrSubjectIsNotMember
		}
	}

	// Список тех, кто может отменять приглашение
	allowedSubjects := []string{
		chat.ChiefID,           // Главный администратор
		invitation.SubjectID,   // Пригласивший
		invitation.RecipientID, // Приглашаемый
	}
	// Проверить, может ли пользователь отменить приглашение
	if !slices.Contains(allowedSubjects, in.SubjectID) {
		return ErrSubjectUserNotAllowed
	}

	// Удаляем приглашение из чата
	if err := chat.RemoveInvitation(in.InvitationID); err != nil {
		return err
	}

	// Сохранить чат в репозиторий
	if err := c.Repo.Upsert(chat); err != nil {
		return err
	}

	return nil
}
