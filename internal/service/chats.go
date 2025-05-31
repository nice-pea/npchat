package service

import (
	"errors"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Chats сервис, объединяющий случаи использования(юзкейсы) в контексте сущности
type Chats struct {
	//ChatsRepo         domain.ChatsRepository
	//MembersRepo       domain.MembersRepository
	ChatAggregateRepo domain.ChatAggregateRepository
}

// WhichParticipateInput входящие параметры
type WhichParticipateInput struct {
	SubjectUserID string
	UserID        string
}

// Validate валидирует значение отдельно каждого параметры
func (in WhichParticipateInput) Validate() error {
	if err := domain.ValidateID(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}
	if err := domain.ValidateID(in.UserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}

	return nil
}

// WhichParticipateOutput результат запроса чатов
type WhichParticipateOutput struct {
	Chats []domain.ChatAggregate
}

// WhichParticipate возвращает список чатов, в которых участвует пользователь
func (c *Chats) WhichParticipate(in WhichParticipateInput) (WhichParticipateOutput, error) {
	// Валидировать параметры
	var err error
	if err = in.Validate(); err != nil {
		return WhichParticipateOutput{}, err
	}

	// Пользователь может запрашивать только свой список чатов
	if in.UserID != in.SubjectUserID {
		return WhichParticipateOutput{}, ErrUnauthorizedChatsView
	}

	// Получить список участников с фильтром по пользователю
	membersFilter := domain.MembersFilter{
		UserID: in.UserID,
	}
	chats, err := c.ChatAggregateRepo.ByParticipantsFilter(membersFilter)
	if err != nil {
		return WhichParticipateOutput{}, err
	}

	return WhichParticipateOutput{Chats: chats}, err
}

// CreateInput входящие параметры
type CreateInput struct {
	Name        string
	ChiefUserID string
}

// CreateOutput результат создания чата
type CreateOutput struct {
	Chat domain.ChatAggregate
}

// Create создает новый чат и участника для главного администратора - пользователя, который создал этот чат
func (c *Chats) Create(in CreateInput) (CreateOutput, error) {
	chat, err := domain.NewChat(in.Name, in.ChiefUserID)
	if err != nil {
		return CreateOutput{}, err
	}

	// Сохранить чат в репозиторий
	if err := c.ChatAggregateRepo.Upsert(chat); err != nil {
		return CreateOutput{}, err
	}

	return CreateOutput{
		Chat: chat,
	}, nil
}

// UpdateNameInput входящие параметры
type UpdateNameInput struct {
	SubjectID string
	ChatID    string
	NewName   string
}

// Validate валидирует значение отдельно каждого параметры
func (in UpdateNameInput) Validate() error {
	if err := domain.ValidateID(in.ChatID); err != nil {
		return errors.Join(err, ErrInvalidChatID)
	}
	if err := domain.ValidateChatName(in.NewName); err != nil {
		return errors.Join(err, ErrInvalidName)
	}
	if err := domain.ValidateID(in.SubjectID); err != nil {
		return errors.Join(err, ErrInvalidSubjectID)
	}

	return nil
}

// UpdateNameOutput результат обновления названия чата
type UpdateNameOutput struct {
	Chat domain.ChatAggregate
}

// UpdateName обновляет название чата.
// Доступно только для главного администратора этого чата
func (c *Chats) UpdateName(in UpdateNameInput) (UpdateNameOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return UpdateNameOutput{}, err
	}

	// Найти чат для обновления
	chat, err := getChatAggregate(c.ChatAggregateRepo, in.ChatID)
	if err != nil {
		return UpdateNameOutput{}, err
	}

	// Проверить доступ пользователя к этому действию
	if in.SubjectID != chat.ChiefID {
		return UpdateNameOutput{}, ErrSubjectUserIsNotChief
	}

	// Перезаписать с новым значением
	if err = chat.UpdateName(in.NewName); err != nil {
		return UpdateNameOutput{}, err
	}
	if err = c.ChatAggregateRepo.Upsert(chat); err != nil {
		return UpdateNameOutput{}, err
	}

	return UpdateNameOutput{
		Chat: chat,
	}, nil
}
