package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// DeleteMember регистрирует обработчик, позволяющий Удалить участника из чата
func DeleteMember(router http2.Router) {
	type requestBody struct {
		UserID string `json:"user_id"`
	}
	router.HandleFunc(
		"DELETE /chats/{chatID}/members",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			var rb requestBody
			if err := http2.DecodeBody(context, &rb); err != nil {
				return nil, err
			}
			input := service.DeleteMemberInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        http2.PathStr(context, "chatID"),
				UserID:        rb.UserID,
			}
			return nil, context.Services().Members().DeleteMember(input)
		})
}
