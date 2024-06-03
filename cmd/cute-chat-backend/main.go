package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saime-0/cute-chat-backend/internal/app"
	"github.com/saime-0/cute-chat-backend/internal/config"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("[Main] Start cute-chat-backend")
	cfg, err := config.Load()
	if err != nil {
		logrus.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err = app.Start(ctx, cfg); err != nil {
			logrus.Fatal(err)
		}
		os.Exit(0)
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	logrus.Infof("[Main] Received signal %s", <-interrupt)
	cancel()
	logrus.Info("[Main] Cancel context")
	time.Sleep(3 * time.Second)
}
