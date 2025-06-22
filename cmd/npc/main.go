package main

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/saime-0/nice-pea-chat/internal/app"
)

func main() {
	slog.Info("main: starting")
	ctx, cancel := context.WithCancel(context.Background())
	go appRun(ctx)
	waitInterrupt(cancel)
}

// appRun запускает приложение и обрабатывает результат
func appRun(ctx context.Context) {
	err := initCliCommand().Run(ctx, os.Args)
	if errors.Is(err, context.Canceled) {
		slog.Info("main: appRun: exit by context canceled")
		os.Exit(0)
	} else if err != nil {
		slog.Error("main: appRun: " + err.Error())
		os.Exit(1)
	}
	slog.Info("main: appRun: finished")
	os.Exit(0)
}

// waitInterrupt отменяет контекст, когда в приложение поступает сигнал syscall.SIGINT или syscall.SIGTERM
func waitInterrupt(cancel context.CancelFunc) {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	slog.Info("main: waitInterrupt: Received signal " + (<-interrupt).String())
	cancel()
	slog.Info("main: waitInterrupt: Context canceled")
	time.Sleep(3 * time.Second)
}

// initCliCommand создает команду, для разбора аргументов командной строки и запуска приложения
func initCliCommand() *cli.Command {
	var cfg app.Config
	return &cli.Command{
		Name: "npchat",
		Action: func(ctx context.Context, command *cli.Command) error {
			return app.Run(ctx, cfg)
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "sqlite-migrations-dir",
				Destination: &cfg.SQLite.MigrationsDir,
				Usage:       "Путь к директории с миграциями SQLite",
				Value:       "./migrations/repository/sqlite",
			},
			&cli.StringFlag{
				Name:        "http-addr",
				Destination: &cfg.HttpAddr,
				Usage:       "Адрес для запуска HTTP сервера",
				Value:       ":8080",
			},
			// Google
			&cli.StringFlag{
				Name:        "oauth-google-client-id",
				Destination: &cfg.OAuthGoogle.ClientID,
				Usage:       "ID клиента OAuth Google",
			},
			&cli.StringFlag{
				Name:        "oauth-google-client-secret",
				Destination: &cfg.OAuthGoogle.ClientSecret,
				Usage:       "Секрет клиента OAuth Google",
			},
			&cli.StringFlag{
				Name:        "oauth-google-callback-url",
				Destination: &cfg.OAuthGoogle.RedirectURL,
				Usage:       "URL для перенаправления после аутентификации OAuth Google",
			},
			// GitHub
			&cli.StringFlag{
				Name:        "oauth-github-client-id",
				Destination: &cfg.OAuthGitHub.ClientID,
				Usage:       "ID клиента OAuth GitHub",
			},
			&cli.StringFlag{
				Name:        "oauth-github-client-secret",
				Destination: &cfg.OAuthGitHub.ClientSecret,
				Usage:       "Секрет клиента OAuth GitHub",
			},
			&cli.StringFlag{
				Name:        "oauth-github-callback-url",
				Destination: &cfg.OAuthGitHub.RedirectURL,
				Usage:       "URL для перенаправления после аутентификации OAuth GitHub",
			},
		},
	}
}
