package http

import (
	"encoding/json"

	"github.com/saime-0/nice-pea-chat/internal/service"
)

// CreateChat создает новый чат
func (c *Controller) CreateChat(context Context) (any, error) {
	var input service.CreateInput
	if err := json.NewDecoder(context.request.Body).Decode(&input); err != nil {
		return nil, err
	}

	result, err := c.chats.Create(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetChats возвращает список чатов пользователя
func (c *Controller) GetChats(context Context) (any, error) {
	input := service.UserChatsInput{
		SubjectUserID: context.subjectID,
		UserID:        context.subjectID,
	}

	chats, err := c.chats.UserChats(input)
	if err != nil {
		return nil, err
	}

	return chats, nil
}

func (c *Controller) Ping(context Context) (any, error) {
	return "pong", nil
}
