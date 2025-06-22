package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
)

// ChatMembersIn входящие параметры
type ChatMembersIn struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
}

// Validate валидирует значение отдельно каждого параметры
func (in ChatMembersIn) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}

	return nil
}

// ChatMembersOut результат запроса чатов
type ChatMembersOut struct {
	Participants []chatt.Participant
}

// ChatMembers возвращает список участников чата
func (c *Chats) ChatMembers(in ChatMembersIn) (ChatMembersOut, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return ChatMembersOut{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return ChatMembersOut{}, err
	}

	// Пользователь должен быть участником чата
	if !chat.HasParticipant(in.SubjectID) {
		return ChatMembersOut{}, ErrSubjectIsNotMember
	}

	return ChatMembersOut{
		Participants: chat.Participants,
	}, nil
}

// LeaveChatIn входящие параметры
type LeaveChatIn struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
}

// Validate валидирует значение отдельно каждого параметры
func (in LeaveChatIn) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}

	return nil
}

// LeaveChat удаляет участника из чата
func (c *Chats) LeaveChat(in LeaveChatIn) error {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return err
	}

	// Удалить пользователя из чата
	if err = chat.RemoveParticipant(in.SubjectID); err != nil {
		return err
	}

	// Сохранить чат в репозиторий
	if err = c.Repo.Upsert(chat); err != nil {
		return err
	}

	return nil
}

// DeleteMemberIn входящие параметры
type DeleteMemberIn struct {
	SubjectID uuid.UUID
	ChatID    uuid.UUID
	UserID    uuid.UUID
}

// Validate валидирует значение отдельно каждого параметры
func (in DeleteMemberIn) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := domain.ValidateID(in.UserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}

	return nil
}

// DeleteMember удаляет участника чата
func (c *Chats) DeleteMember(in DeleteMemberIn) error {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return err
	}

	// Проверить попытку удалить самого себя
	if in.UserID == in.SubjectID {
		return ErrMemberCannotDeleteHimself
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return err
	}

	// Subject должен быть главным администратором
	if chat.ChiefID != in.SubjectID {
		return ErrSubjectUserIsNotChief
	}

	// Удалить пользователя из чата
	if err = chat.RemoveParticipant(in.UserID); err != nil {
		return err
	}

	// Сохранить чат в репозиторий
	if err = c.Repo.Upsert(chat); err != nil {
		return err
	}

	return nil
}
