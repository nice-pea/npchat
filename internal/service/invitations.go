package service

import (
	"errors"

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
	panic("not implemented")
}

type ReceivedInvitationsIn struct {
	SubjectID uuid.UUID
}

func (in ReceivedInvitationsIn) Validate() error {
	return nil
}

// ReceivedInvitationsOut входящие параметры
type ReceivedInvitationsOut struct {
	// ChatsInvitations карта приглашений, где ключ - chatID, значение - приглашение
	ChatsInvitations map[uuid.UUID]chatt.Invitation
}

// ReceivedInvitations возвращает список приглашений конкретного пользователя в чаты
func (c *Chats) ReceivedInvitations(in ReceivedInvitationsIn) (ReceivedInvitationsOut, error) {
	panic("not implemented")
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
