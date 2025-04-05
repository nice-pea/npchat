package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

// Chats сервис объединяющий случаи использования(юзкейсы) в контексте сущности
type Chats struct {
	ChatsRepo   domain.ChatsRepository
	MembersRepo domain.MembersRepository
}

// UserChatsInput входящие параметры
type UserChatsInput struct {
	SubjectUserID string
	UserID        string
}

// Validate валидирует значение отдельно каждого параметры
func (in UserChatsInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrInvalidSubjectUserID)
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return errors.Join(err, ErrInvalidUserID)
	}

	return nil
}

// UserChats возвращает список чатов в которых участвует пользователь
func (c *Chats) UserChats(in UserChatsInput) ([]domain.Chat, error) {
	// Валидировать параметры
	var err error
	if err = in.Validate(); err != nil {
		return nil, err
	}

	// Пользователь может запрашивать только свой список чатов
	if in.UserID != in.SubjectUserID {
		return nil, ErrCannotViewSomeoneElseChats
	}

	// Получить список участников с фильтром по пользователю
	members, err := c.MembersRepo.List(domain.MembersFilter{
		UserID: in.UserID,
	})
	if err != nil {
		return nil, err
	}

	// Если нет участников, то запрашивать чаты не надо
	if len(members) == 0 {
		return nil, nil
	}

	// Собрать ID чатов к которым принадлежат участники
	chatIds := make([]string, len(members))
	for i, member := range members {
		chatIds[i] = member.ChatID
	}

	// Вернуть список чатов с фильтром по ID
	chats, err := c.ChatsRepo.List(domain.ChatsFilter{
		IDs: chatIds,
	})

	return chats, err
}

// CreateInput входящие параметры
type CreateInput struct {
	Name        string
	ChiefUserID string
}

// Validate валидирует значение отдельно каждого параметры
func (in CreateInput) Validate() error {
	chat := domain.Chat{
		Name:        in.Name,
		ChiefUserID: in.ChiefUserID,
	}
	if err := chat.ValidateName(); err != nil {
		return errors.Join(err, ErrInvalidName)
	}
	if err := chat.ValidateChiefUserID(); err != nil {
		return errors.Join(err, ErrInvalidChiefUserID)
	}

	return nil
}

// CreateOutput результат создания чата
type CreateOutput struct {
	Chat        domain.Chat
	ChiefMember domain.Member
}

// Create создает новый чат и участника для главного администратора - пользователя который создал этот чат
func (c *Chats) Create(in CreateInput) (CreateOutput, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return CreateOutput{}, err
	}

	// Сохранить чат в репозиторий
	newChat := domain.Chat{
		ID:          uuid.NewString(),
		Name:        in.Name,
		ChiefUserID: in.ChiefUserID,
	}
	if err := c.ChatsRepo.Save(newChat); err != nil {
		return CreateOutput{}, err
	}

	// Создать участника
	member := domain.Member{
		ID:     uuid.NewString(),
		UserID: newChat.ChiefUserID,
		ChatID: newChat.ID,
	}
	if err := c.MembersRepo.Save(member); err != nil {
		return CreateOutput{}, err
	}

	return CreateOutput{
		Chat:        newChat,
		ChiefMember: member,
	}, nil
}

// UpdateNameInput входящие параметры
type UpdateNameInput struct {
	SubjectUserID string
	ChatID        string
	NewName       string
}

// Validate валидирует значение отдельно каждого параметры
func (in UpdateNameInput) Validate() error {
	chat := domain.Chat{
		ID:          in.ChatID,
		Name:        in.NewName,
		ChiefUserID: in.SubjectUserID,
	}
	if err := chat.ValidateID(); err != nil {
		return errors.Join(err, ErrInvalidID)
	}
	if err := chat.ValidateName(); err != nil {
		return errors.Join(err, ErrInvalidName)
	}
	if err := chat.ValidateChiefUserID(); err != nil {
		return errors.Join(err, ErrInvalidChiefUserID)
	}

	return nil
}

// UpdateName обновляет название чата.
// Доступно только для главного администратора этого чата
func (c *Chats) UpdateName(in UpdateNameInput) (domain.Chat, error) {
	// Валидировать параметры
	if err := in.Validate(); err != nil {
		return domain.Chat{}, err
	}

	// Найти чат для обновления
	chat, err := getChat(c.ChatsRepo, in.ChatID)
	if err != nil {
		return domain.Chat{}, err
	}

	// Проверить доступ пользователя к этому действию
	if in.SubjectUserID != chat.ChiefUserID {
		return domain.Chat{}, ErrSubjectUserIsNotChief
	}

	// Перезаписать с новым значением
	updatedChat := domain.Chat{
		ID:          chat.ID,
		Name:        in.NewName,
		ChiefUserID: chat.ChiefUserID,
	}
	if err = c.ChatsRepo.Save(updatedChat); err != nil {
		return domain.Chat{}, err
	}

	return updatedChat, nil
}
