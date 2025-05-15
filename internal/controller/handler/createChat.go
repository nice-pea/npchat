package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

type requestBody struct {
	Name        string `json:"name"`
	ChiefUserID string `json:"chief_user_id"`
}

func CreateChat(router http2.Router) {
	router.HandleFunc("POST /chats", func(context http2.Context) (any, error) {
		var rb requestBody
		if err := http2.DecodeBody(context, &rb); err != nil {
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
	})
}
