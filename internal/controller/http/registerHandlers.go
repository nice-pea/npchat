package http

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
