package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saime-0/nice-pea-chat/internal/app"
	"github.com/saime-0/nice-pea-chat/internal/config"
)

func main() {
	log.Println("[Main] Start cute-chat-backend")
	cfg, err := config.Load()
	if err != nil {
		log.Println(err)
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
	log.Printf("[Main] Received signal %s", <-interrupt)
	cancel()
	log.Println("[Main] Cancel context")
	time.Sleep(3 * time.Second)
}
