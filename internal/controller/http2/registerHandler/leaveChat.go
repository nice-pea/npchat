package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Покинуть чат
func LeaveChat(router http2.Router) {
	router.HandleFunc(
		"POST /chats/{chatID}/leave",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			input := service.LeaveChatInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        http2.PathStr(context, "chatID"),
			}
			return nil, context.Services().Members().LeaveChat(input)
		})
}
