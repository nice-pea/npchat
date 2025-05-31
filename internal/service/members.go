package service

import (
	"errors"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
)

// ChatMembersInput входящие параметры
type ChatMembersInput struct {
	SubjectID string
	ChatID    string
}

// Validate валидирует значение отдельно каждого параметры
func (in ChatMembersInput) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}

	return nil
}

// ChatMembersOutput результат запроса чатов
type ChatMembersOutput struct {
	Participants []chatt.Participant
}

// ChatMembers возвращает список участников чата
func (c *Chats) ChatMembers(in ChatMembersInput) (ChatMembersOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return ChatMembersOutput{}, err
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return ChatMembersOutput{}, err
	}

	// Пользователь должен быть участником чата
	if !chat.HasParticipant(in.SubjectID) {
		return ChatMembersOutput{}, ErrSubjectIsNotMember
	}

	return ChatMembersOutput{
		Participants: chat.Participants,
	}, nil
}

// LeaveChatInput входящие параметры
type LeaveChatInput struct {
	SubjectID string
	ChatID    string
}

// Validate валидирует значение отдельно каждого параметры
func (in LeaveChatInput) Validate() error {
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}

	return nil
}

// LeaveChat удаляет участника из чата
func (c *Chats) LeaveChat(in LeaveChatInput) error {
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

// DeleteMemberInput входящие параметры
type DeleteMemberInput struct {
	SubjectUserID string
	ChatID        string
	UserID        string
}

// Validate валидирует значение отдельно каждого параметры
func (in DeleteMemberInput) Validate() error {
	if err := domain.ValidateID(in.SubjectUserID); err != nil {
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
func (c *Chats) DeleteMember(in DeleteMemberInput) error {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return err
	}

	// Проверить попытку удалить самого себя
	if in.UserID == in.SubjectUserID {
		return ErrMemberCannotDeleteHimself
	}

	// Найти чат
	chat, err := chatt.Find(c.Repo, chatt.Filter{ID: in.ChatID})
	if err != nil {
		return err
	}

	// Subject должен быть главным администратором
	if chat.ChiefID != in.SubjectUserID {
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
