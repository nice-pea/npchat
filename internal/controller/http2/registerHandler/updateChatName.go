package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Обновить название чата
func UpdateChatName(router http2.Router) {
	type requestBody struct {
		NewName string `json:"new_name"`
	}
	router.HandleFunc(
		"PUT /chats/{chatID}/name",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}

			input := service.UpdateNameInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        http2.PathStr(context, "chatID"),
				NewName:       rb.NewName,
			}

			return context.Services().Chats().UpdateName(input)
		})
}
