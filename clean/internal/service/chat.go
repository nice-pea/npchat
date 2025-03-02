package service

import (
	"errors"

	"github.com/google/uuid"

	"github.com/saime-0/nice-pea-chat/internal/common"
	"github.com/saime-0/nice-pea-chat/internal/domain"
)

type Chats struct {
	ChatsRepo   domain.ChatsRepository
	MembersRepo domain.MembersRepository
	History     History
}

const LimitCreatedChats = 1000

type ChatsCreateIn struct {
	Name    string
	OwnerID string
}
type ChatCreateErr struct {
	Error error
}
type ChatWasCreated struct {
	Chat  domain.Chat
	Owner domain.Member
}

var ErrCreateChatLimitExceeded = errors.New("create chat limit exceeded")

func (c *Chats) Create(in ChatsCreateIn) (*domain.Chat, error) {
	// Проверить лимит созданных чатов
	members, err := c.MembersRepo.List(domain.MembersFilter{
		UserID:  in.OwnerID,
		IsOwner: common.Pointer(true),
	})
	if err != nil {
		c.History.Write(ChatCreateErr{Error: err})
		return nil, err
	}
	if len(members) >= LimitCreatedChats {
		c.History.Write(ChatCreateErr{Error: ErrCreateChatLimitExceeded})
		return nil, ErrCreateChatLimitExceeded
	}

	// Создать чат
	chat := domain.Chat{
		ID:   uuid.NewString(),
		Name: in.Name,
	}
	if err = c.ChatsRepo.Save(chat); err != nil {
		c.History.Write(ChatCreateErr{Error: err})
		return nil, err
	}

	// Добавить участника
	member := domain.Member{
		ID:      uuid.NewString(),
		UserID:  in.OwnerID,
		ChatID:  chat.ID,
		IsOwner: true,
	}
	err = c.MembersRepo.Save(member)
	if err != nil {
		c.History.Write(ChatCreateErr{Error: err})
		return nil, err
	}

	// Запись в историю
	c.History.Write(ChatWasCreated{
		Chat:  chat,
		Owner: member,
	})

	return &chat, nil
}

func (c *Chats) Delete(memberID uint) error {
	return nil
}

type ChatsFilter struct {
	ID      uint
	UserIDs []uint
}

func (c *Chats) List(memberID uint) error {
	return nil
}

func (c *Chats) Members(chatID uint) ([]domain.Member, error) {
	return nil, nil
}
