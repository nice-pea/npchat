package http

func registerHandlers(c *Controller) {
	c.mux.HandleFunc("POST /chats", c.modulation(c.CreateChat))
	c.mux.HandleFunc("GET /chats", c.modulation(c.GetChats))
	c.mux.HandleFunc("/ping", c.modulation(c.Ping))
}
