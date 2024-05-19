package main

import (
	"context"
	"github.com/saime-0/cute-chat-backend/internal/app"
	"github.com/saime-0/cute-chat-backend/internal/config"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())

	app.Start(ctx, wg, cfg)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-quit
	go func() {
		time.Sleep(3 * time.Second)
		cancel()
	}()
	wg.Wait()
}
