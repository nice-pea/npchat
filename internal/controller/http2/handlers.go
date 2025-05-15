package http2

import (
	"encoding/json"

	"github.com/saime-0/nice-pea-chat/internal/service"
)

func (c *Controller) Ping(context Context) (any, error) {
	return "pong", nil
}

func (c *Controller) LoginByPassword(context Context) (any, error) {
	var input service.AuthnPasswordLoginInput
	if err := json.NewDecoder(context.request.Body).Decode(&input); err != nil {
		return nil, err
	}

	session, err := c.authnPassword.Login(input)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// MyChats возвращает список чатов пользователя
func (c *Controller) MyChats(context Context) (any, error) {
	input := service.UserChatsInput{
		SubjectUserID: context.session.UserID,
		UserID:        context.session.UserID,
	}

	chats, err := c.chats.UserChats(input)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

// CreateChat создает новый чат
func (c *Controller) CreateChat(context Context) (any, error) {
	var rb struct {
		Name        string `json:"name"`
		ChiefUserID string `json:"chief_user_id"`
	}
	if err := json.NewDecoder(context.request.Body).Decode(&rb); err != nil {
		return nil, err
	}
	input := service.CreateInput{
		Name:        rb.Name,
		ChiefUserID: rb.ChiefUserID,
	}

	result, err := c.chats.Create(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Controller) UpdateChatName(context Context) (any, error) {
	return nil, nil
}

func (c *Controller) LeaveChat(context Context) (any, error) {
	return nil, nil
}

func (c *Controller) ChatMembers(context Context) (any, error) {
	input := service.ChatMembersInput{
		SubjectUserID: context.session.UserID,
		ChatID:        context.request.PathValue("chatID"),
	}

	chats, err := c.members.ChatMembers(input)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (c *Controller) MyInvitations(context Context) (any, error) {
	return nil, nil
}

func (c *Controller) ChatInvitations(context Context) (any, error) {
	return nil, nil
}

func (c *Controller) SendInvitation(context Context) (any, error) {
	return nil, nil
}

func (c *Controller) AcceptInvitation(context Context) (any, error) {
	return nil, nil
}

func (c *Controller) CancelInvitation(context Context) (any, error) {
	return nil, nil
}

func (c *Controller) DeleteMember(context Context) (any, error) {
	return nil, nil
}
