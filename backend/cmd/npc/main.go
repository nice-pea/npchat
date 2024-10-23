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
	log.Println("[main] Start nice-pea-chat")
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
	log.Printf("[main] Received signal %s", <-interrupt)
	cancel()
	log.Println("[main] Cancel context")
	time.Sleep(3 * time.Second)
}
