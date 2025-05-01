package http

func (c *Controller) registerHandlers() {
	clientChain := []middleware{
		c.requireRequestID,
		requireAcceptJson,
		c.requireAuthorizedSession,
	}
	c.HandleFunc("POST /chats", c.CreateChat, clientChain...)
	c.HandleFunc("GET /chats", c.GetChats, clientChain...)
	c.HandleFunc("/ping", c.Ping)
}
