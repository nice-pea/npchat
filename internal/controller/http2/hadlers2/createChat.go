package hadlers2

import (
	"github.com/saime-0/nice-pea-chat/internal/controller"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// CreateChat обработчик для создания нового чата
type CreateChat struct{}

func (c *CreateChat) HandlerParams() controller.HandlerParams {
	return controller.HandlerParams{
		Method: "POST",
		Path:   "/chats",
	}
}

func (c *CreateChat) HandlerFunc(context controller.Context) (any, error) {
	var rb struct {
		Name        string `json:"name"`
		ChiefUserID string `json:"chief_user_id"`
	}
	if err := controller.DecodeBody(context, &rb); err != nil {
		return nil, err
	}

	input := service.CreateInput{
		Name:        rb.Name,
		ChiefUserID: rb.ChiefUserID,
	}

	result, err := context.Services().Chats().Create(input)
	if err != nil {
		return nil, err
	}

	return result, nil
}
