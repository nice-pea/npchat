package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"

	"github.com/saime-0/nice-pea-chat/internal/config"
	"github.com/saime-0/nice-pea-chat/internal/http/handlers"
	"github.com/saime-0/nice-pea-chat/internal/httpserver"
	"github.com/saime-0/nice-pea-chat/internal/repository/postgres"
)

func Start(ctx context.Context, cfg *config.Config) error {
	var wg sync.WaitGroup
	db, err := pgx.Connect(ctx, cfg.DB)
	if err != nil {
		return fmt.Errorf("pgx.Connect: %w", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		db.Close(context.Background())
	}()

	commonRepository := postgres.NewCommonRepository(db)

	httpServer := httpserver.New(
		ctx, cfg.Listen,
		httpserver.Handlers{
			&handlers.Healthcheck{},
			&handlers.Auth{},
			&handlers.UserByToken{},
			&handlers.UserByID{},
			&handlers.UserUpdate{},
			&handlers.UserChats{},
		},
	)
	wg.Add(1)
	go func() {
		<-httpServer.Done()
		wg.Done()
	}()

	select {
	case <-ctx.Done():
		log.Println("[App] Receive ctx.Done, wait when components stop the work")
		wg.Wait()
		log.Println("[App] Components done the work")
		return nil
	case err = <-httpServer.Notify():
		return fmt.Errorf("received notify from httpServer: %w", err)
	}
}
