package http2

import "github.com/saime-0/nice-pea-chat/internal/controller"

//func (c *Controller) RegisterHandlers(handlers ...controller.Handler2) error {
//	for i := range handlers {
//		p := handlers[i].HandlerParams()
//		handlerFunc := handlers[i].HandlerFunc
//		c.HandleFunc(p.Method+" "+p.Path, handlerFunc)
//	}
//
//	return nil
//}

func (c *Controller) registerHandlers() {
	c.HandleFunc("/ping", c.Ping)

	// Публичные эндпоинты, без аутентификации
	clientPubChain := []middleware{
		c.requireRequestID,
		requireAcceptJson,
	}
	c.HandleFunc("POST /login/password", c.LoginByPassword, clientPubChain...)

	// Эндпоинты, требующие аутентификации
	clientAuthChain := []middleware{
		c.requireRequestID,
		requireAcceptJson,
		c.requireAuthorizedSession,
	}
	// Чат
	c.HandleFunc("GET /chats", c.MyChats, clientAuthChain...)
	c.HandleFunc("POST /chats", c.CreateChat, clientAuthChain...)
	c.HandleFunc("PUT /chats/{chatID}/name", c.UpdateChatName, clientAuthChain...)

	// Участники
	c.HandleFunc("POST /chats/{chatID}/leave", c.LeaveChat, clientAuthChain...)
	c.HandleFunc("GET /chats/{chatID}/members", c.ChatMembers, clientAuthChain...)
	c.HandleFunc("DELETE /chats/{chatID}/members/{memberID}", c.DeleteMember, clientAuthChain...)

	// Приглашениями
	c.HandleFunc("GET /invitations", c.MyInvitations, clientAuthChain...)
	c.HandleFunc("GET /chats/{chatID}/invitations", c.ChatInvitations, clientAuthChain...)
	c.HandleFunc("POST /invitations", c.SendInvitation, clientAuthChain...)
	c.HandleFunc("POST /invitations/{invitationID}/accept", c.AcceptInvitation, clientAuthChain...)
	c.HandleFunc("POST /invitations/{invitationID}/cancel", c.CancelInvitation, clientAuthChain...)

}
