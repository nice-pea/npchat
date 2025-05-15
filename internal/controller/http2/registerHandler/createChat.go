package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Создать чат
func CreateChat(router http2.Router) {
	type requestBody struct {
		Name string `json:"name"`
	}
	router.HandleFunc(
		"POST /chats",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.CreateInput{
				ChiefUserID: context.Session().UserID,
				Name:        rb.Name,
			}

			result, err := context.Services().Chats().Create(input)
			if err != nil {
				return nil, err
			}

			return result, nil
		})
}
