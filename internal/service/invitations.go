package service

import (
	"errors"
	"slices"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
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
	return nil
}

// AcceptInvitation добавляет пользователя в чат, путем принятия приглашения
func (c *Chats) AcceptInvitation(in AcceptInvitationIn) error {
	return nil
}

type CancelInvitationIn struct {
	SubjectID    uuid.UUID
	InvitationID uuid.UUID
}

func (in CancelInvitationIn) Validate() error {
	return nil
}

// CancelInvitation отменяет приглашение
func (c *Chats) CancelInvitation(in CancelInvitationIn) error {
	return nil
}
