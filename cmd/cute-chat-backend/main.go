package main

import (
	"github.com/saime-0/cute-chat-backend/internal/app"
	"github.com/saime-0/cute-chat-backend/internal/config"
)

func main() {
	app.Start(config.Load())
}
