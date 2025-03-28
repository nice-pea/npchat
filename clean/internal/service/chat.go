package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type Chats struct {
	ChatsRepo   domain.ChatsRepository
	MembersRepo domain.MembersRepository
	History     History
}

type ChatsWhereUserIsMemberInput struct {
	SubjectUserID string
	UserID        string
}

func (in ChatsWhereUserIsMemberInput) validate() error {
	if in.SubjectUserID == "" {
		return errors.New("subjectUserID is empty")
	}
	if in.UserID == "" {
		return errors.New("UserID is empty")
	}
	if in.UserID != in.SubjectUserID {
		return errors.New("должны быть одинаковыми")
	}

	return nil
}

// ChatsWhereUserIsMember возвращает список чатов в которых участвует пользователь
func (c *Chats) ChatsWhereUserIsMember(in ChatsWhereUserIsMemberInput) ([]domain.Chat, error) {
	// Валидировать параметры
	var err error
	if err = in.validate(); err != nil {
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

type CreateInput struct {
	Name        string
	ChiefUserID string
}
type CreateOutput struct {
	Chat        domain.Chat
	ChiefMember domain.Member
}

// Create создает новый чат и участника для главного администратора - пользователя который создал этот чат
func (c *Chats) Create(in CreateInput) (CreateOutput, error) {
	newChat := domain.Chat{
		ID:          uuid.NewString(),
		Name:        in.Name,
		ChiefUserID: in.ChiefUserID,
	}

	// Валидация создаваемого чата
	if err := newChat.ValidateName(); err != nil {
		return CreateOutput{}, err
	}
	if err := newChat.ValidateChiefUserID(); err != nil {
		return CreateOutput{}, err
	}

	// Сохранить чат в репозиторий
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
