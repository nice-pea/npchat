package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/saime-0/nice-pea-chat/internal/app"
)

func main() {
	slog.Info("main: Starting")
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := app.Run(ctx); err != nil {
			slog.Error("app.Run:" + err.Error())
			os.Exit(1)
		}
		os.Exit(0)
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("main: Received signal " + (<-interrupt).String())
	cancel()
	slog.Info("main: Context canceled")
	time.Sleep(3 * time.Second)
}
