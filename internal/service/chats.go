package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/nice-pea/npchat/internal/domain"
	"github.com/nice-pea/npchat/internal/domain/chatt"
)

// Chats сервис, объединяющий случаи использования(юзкейсы) в контексте агрегата чатов
type Chats struct {
	Repo chatt.Repository
}

// WhichParticipateIn входящие параметры
type WhichParticipateIn struct {
	SubjectID uuid.UUID
	UserID    uuid.UUID // TODO: удалить
}

// Validate валидирует значение отдельно каждого параметры
func (in WhichParticipateIn) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.UserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}

	return nil
}

// WhichParticipateOut результат запроса чатов
type WhichParticipateOut struct {
	Chats []chatt.Chat
}

// WhichParticipate возвращает список чатов, в которых участвует пользователь
func (c *Chats) WhichParticipate(in WhichParticipateIn) (WhichParticipateOut, error) {
	// Валидировать параметры
	var err error
	if err = in.Validate(); err != nil {
		return WhichParticipateOut{}, err
	}

	// Пользователь может запрашивать только свой список чатов
	if in.UserID != in.SubjectID {
		return WhichParticipateOut{}, ErrUnauthorizedChatsView
	}

	// Получить список участников с фильтром по пользователю
	chats, err := c.Repo.List(chatt.Filter{
		ParticipantID: in.UserID,
	})
	if err != nil {
		return WhichParticipateOut{}, err
	}

	return WhichParticipateOut{Chats: chats}, err
}


// UpdateNameIn входящие параметры
type UpdateNameIn struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
	NewName   string
}

// Validate валидирует значение отдельно каждого параметры
func (in UpdateNameIn) Validate() error {
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := chatt.ValidateChatName(in.NewName); err != nil {
		return errors.Join(err, ErrInvalidName)
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// UpdateNameOut результат обновления названия чата
type UpdateNameOut struct {
	Chat chatt.Chat
}

// UpdateName обновляет название чата.
// Доступно только для главного администратора этого чата
func (c *Chats) UpdateName(in UpdateNameIn) (UpdateNameOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return UpdateNameOut{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return UpdateNameOut{}, err
	}

	// Проверить доступ пользователя к этому действию
	if in.SubjectID != chat.ChiefID {
		return UpdateNameOut{}, ErrSubjectUserIsNotChief
	}

	// Перезаписать с новым значением
	if err = chat.UpdateName(in.NewName); err != nil {
		return UpdateNameOut{}, err
	}
	if err = c.Repo.Upsert(chat); err != nil {
		return UpdateNameOut{}, err
	}

	return UpdateNameOut{
		Chat: chat,
	}, nil
}

// AcceptInvitation
// CancelInvitation
// ChatInvitations
// ChatMembers
// CreateChat
// DeleteMember
// LeaveChat
// ReceivedInvitations
// SendInvitation
// UpdateName
// WhichParticipate

// accept_invitation
// cancel_invitation
// chat_invitations
// chat_members
// create_chat
// delete_member
// leave_chat
// received_invitations
// send_invitation
// update_name
// which_participate

// accept_invitation cancel_invitation chat_invitations chat_members create_chat delete_member leave_chat received_invitations send_invitation update_name which_participate 
