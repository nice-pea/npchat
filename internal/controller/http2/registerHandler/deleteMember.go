package registerHandler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/http2/middleware"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

// Удалить участника из чата
func DeleteMember(router http2.Router) {
	router.HandleFunc(
		"DELETE /chats/{chatID}/members/{memberID}",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			input := service.DeleteMemberInput{
				SubjectUserID: context.Session().UserID,
				ChatID:        http2.PathStr(context, "chatID"),
				UserID:        "",
			}
			return nil, context.Services().Members().DeleteMember()
		})
}
