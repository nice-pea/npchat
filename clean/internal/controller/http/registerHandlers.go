package http

func registerHandlers(c *Controller) {
	clientChain := []middleware{
		requireRequestID,
		requireAcceptJson,
		requireAuthorizedSession,
	}
	c.HandleFunc("POST /chats", c.CreateChat, clientChain...)
	c.HandleFunc("GET /chats", c.GetChats, clientChain...)
	c.HandleFunc("/ping", c.Ping)
}
