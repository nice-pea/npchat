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
	History     History
}

// UserChatsInput параметры для запроса чатов в которых участвует пользователь
type UserChatsInput struct {
	SubjectUserID string
	UserID        string
}

var (
	ErrUserChatsInputSubjectUserIDValidate = errors.New("некорректный SubjectUserID")
	ErrUserChatsInputUserIDValidate        = errors.New("некорректный UserID")
	ErrUserChatsInputEqualUserIDsValidate  = errors.New("UserID и SubjectUserID не совпадают")
)

// Validate валидирует параметры для запроса чатов в которых участвует пользователь
func (in UserChatsInput) Validate() error {
	if err := uuid.Validate(in.SubjectUserID); err != nil {
		return errors.Join(err, ErrUserChatsInputSubjectUserIDValidate)
	}
	if err := uuid.Validate(in.UserID); err != nil {
		return errors.Join(err, ErrUserChatsInputUserIDValidate)
	}
	if in.UserID != in.SubjectUserID {
		return ErrUserChatsInputEqualUserIDsValidate
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

	// Получить список участников с фильтром по пользователю
	var members []domain.Member
	if members, err = c.MembersRepo.List(domain.MembersFilter{
		UserID: in.UserID,
	}); err != nil {
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
	return c.ChatsRepo.List(domain.ChatsFilter{
		IDs: chatIds,
	})
}

// CreateInput входящие параметры для создания чата
type CreateInput struct {
	Name        string
	ChiefUserID string
}

var (
	ErrCreateInputNameValidate        = errors.New("некорректный Name")
	ErrCreateInputChiefUserIDValidate = errors.New("некорректный ChiefUserID")
)

// Validate валидирует параметры для создания чата
func (in CreateInput) Validate() error {
	chat := domain.Chat{
		Name:        in.Name,
		ChiefUserID: in.ChiefUserID,
	}
	if err := chat.ValidateName(); err != nil {
		return errors.Join(err, ErrCreateInputNameValidate)
	}
	if err := chat.ValidateChiefUserID(); err != nil {
		return errors.Join(err, ErrCreateInputChiefUserIDValidate)
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

// UpdateNameInput входящие параметры для обновления названия чата
type UpdateNameInput struct {
	SubjectUserID string
	ChatID        string
	NewName       string
}

var (
	ErrUpdateNameIDValidate          = errors.New("некорректный ID")
	ErrUpdateNameNameValidate        = errors.New("некорректный Name")
	ErrUpdateNameChiefUserIDValidate = errors.New("некорректный ChiefUserID")
)

// Validate валидирует параметры для обновления названия чата
func (in UpdateNameInput) Validate() error {
	chat := domain.Chat{
		ID:          in.ChatID,
		Name:        in.NewName,
		ChiefUserID: in.SubjectUserID,
	}
	if err := chat.ValidateID(); err != nil {
		return errors.Join(err, ErrUpdateNameIDValidate)
	}
	if err := chat.ValidateName(); err != nil {
		return errors.Join(err, ErrUpdateNameNameValidate)
	}
	if err := chat.ValidateChiefUserID(); err != nil {
		return errors.Join(err, ErrUpdateNameChiefUserIDValidate)
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
	chats, err := c.ChatsRepo.List(domain.ChatsFilter{
		IDs: []string{in.ChatID},
	})
	if err != nil {
		return domain.Chat{}, err
	}
	if len(chats) != 1 {
		return domain.Chat{}, errors.New("чат не найден")
	}

	// Проверить доступ пользователя к этому действию
	if in.SubjectUserID != chats[0].ChiefUserID {
		return domain.Chat{}, errors.New("доступно только главному администратору")
	}

	// Перезаписать с новым значением
	updatedChat := domain.Chat{
		ID:          chats[0].ID,
		Name:        in.NewName,
		ChiefUserID: chats[0].ChiefUserID,
	}
	if err = c.ChatsRepo.Save(updatedChat); err != nil {
		return domain.Chat{}, err
	}

	return updatedChat, nil
}
