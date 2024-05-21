package main

import (
	"context"
	"github.com/saime-0/cute-chat-backend/internal/app"
	"github.com/saime-0/cute-chat-backend/internal/config"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err = app.Start(ctx, cfg); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Printf("main - signal: " + (<-interrupt).String())
	cancel()
	time.Sleep(3 * time.Second)
}
