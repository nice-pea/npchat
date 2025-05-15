package handler

import (
	"github.com/saime-0/nice-pea-chat/internal/controller/http2"
	"github.com/saime-0/nice-pea-chat/internal/controller/middleware"
)

// Удалить участника из чата
func RegisterDeleteMemberHandler(router http2.Router) {
	router.HandleFunc(
		"DELETE /chats/{chatID}/members/{memberID}",
		middleware.ClientAuthChain,
		func(context http2.Context) (any, error) {
			return "not implemented", nil
		})
}
