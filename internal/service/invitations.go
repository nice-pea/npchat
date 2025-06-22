package service

import (
	"errors"
	"slices"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
)

// ChatInvitationsIn параметры для запроса приглашений конкретного чата
type ChatInvitationsIn struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
}

// Validate валидирует параметры для запроса приглашений конкретного чата
func (in ChatInvitationsIn) Validate() error {
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// ChatInvitationsOut результат запроса приглашений конкретного чата
type ChatInvitationsOut struct {
	Invitations []chatt.Invitation
}

// ChatInvitations возвращает список приглашений в конкретный чат.
// Если SubjectID является администратором, то возвращается все приглашения в данный чат,
// иначе только те приглашения, которые отправил именно пользователь.
func (c *Chats) ChatInvitations(in ChatInvitationsIn) (ChatInvitationsOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return ChatInvitationsOut{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return ChatInvitationsOut{}, err
	}

	// Проверить является ли пользователь участником чата
	if !chat.HasParticipant(in.SubjectID) {
		return ChatInvitationsOut{}, ErrSubjectIsNotMember
	}

	// Сохранить сначала все приглашения
	invitations := chat.Invitations

	// Если пользователь не является администратором,
	// то оставить только те приглашения, которые отправил именно пользователь.
	if chat.ChiefID != in.SubjectID {
		invitations = chat.SubjectInvitations(in.SubjectID)
	}

	return ChatInvitationsOut{
		Invitations: invitations,
	}, err
}

type ReceivedInvitationsIn struct {
	SubjectID uuid.UUID
}

func (in ReceivedInvitationsIn) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// ReceivedInvitationsOut входящие параметры
type ReceivedInvitationsOut struct {
	// ChatsInvitations карта приглашений, где ключ - chatID, значение - приглашение
	ChatsInvitations map[uuid.UUID]chatt.Invitation
}

// ReceivedInvitations возвращает список приглашений конкретного пользователя в чаты
func (c *Chats) ReceivedInvitations(in ReceivedInvitationsIn) (ReceivedInvitationsOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return ReceivedInvitationsOut{}, err
	}

	// Найти чат
	chats, err := c.Repo.List(chatt.Filter{
		InvitationRecipientID: in.SubjectID,
	})
	if err != nil {
		return ReceivedInvitationsOut{}, err
	}

	// Если нет чатов, вернут пустой список
	if len(chats) == 0 {
		return ReceivedInvitationsOut{}, nil
	}

	// Собрать приглашения, полученные пользователем
	invitations := make(map[uuid.UUID]chatt.Invitation, len(chats))
	for _, chat := range chats {
		invitations[chat.ID], _ = chat.RecipientInvitation(in.SubjectID)
	}

	return ReceivedInvitationsOut{
		ChatsInvitations: invitations,
	}, err
}

type SendInvitationIn struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
}

func (in SendInvitationIn) Validate() error {
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

type SendInvitationOut struct {
	Invitation chatt.Invitation
}

// SendInvitation отправляет приглашения пользователю от участника чата
func (c *Chats) SendInvitation(in SendInvitationIn) (SendInvitationOut, error) {
	if err := in.Validate(); err != nil {
		return SendInvitationOut{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return SendInvitationOut{}, err
	}

	// Создать приглашение
	inv, err := chatt.NewInvitation(in.SubjectID, in.UserID)
	if err != nil {
		return SendInvitationOut{}, err
	}

	// Добавить приглашение в чат
	if err = chat.AddInvitation(inv); err != nil {
		return SendInvitationOut{}, err
	}

	// Сохранить чат в репозиторий
	if err = c.Repo.Upsert(chat); err != nil {
		return SendInvitationOut{}, err
	}

	return SendInvitationOut{
		Invitation: inv,
	}, nil
}

type AcceptInvitationIn struct {
	SubjectID    uuid.UUID
	InvitationID uuid.UUID
}

func (in AcceptInvitationIn) Validate() error {
	if err := domain.ValidateID(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}

	return nil
}

// AcceptInvitation добавляет пользователя в чат, путем принятия приглашения
func (c *Chats) AcceptInvitation(in AcceptInvitationIn) error {
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

type CancelInvitationIn struct {
	SubjectID    uuid.UUID
	InvitationID uuid.UUID
}

func (in CancelInvitationIn) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return ErrInvalidSubjectID
	}
	if err := domain.ValidateID(in.InvitationID); err != nil {
		return ErrInvalidInvitationID
	}

	return nil
}

// CancelInvitation отменяет приглашение
func (c *Chats) CancelInvitation(in CancelInvitationIn) error {
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
	allowedSubjects := []uuid.UUID{
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
