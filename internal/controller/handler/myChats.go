package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Получить список чатов пользователя
func RegisterMyChatsHandler(router http2.Router) {
	router.HandleFunc(
		"GET /chats",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			input := service.UserChatsInput{
				SubjectUserID: context.Session().UserID,
				UserID:        context.Session().UserID,
			}

			chats, err := context.Services().Chats().UserChats(input)
			if err != nil {
				return nil, err
			}

			return chats, nil
		})
}
