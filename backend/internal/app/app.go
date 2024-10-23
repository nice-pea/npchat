package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/config"
	"github.com/saime-0/nice-pea-chat/internal/http/handlers"
	"github.com/saime-0/nice-pea-chat/internal/httpserver"
	"github.com/saime-0/nice-pea-chat/internal/repository/postgres"
)

func Start(ctx context.Context, cfg config.Config) error {
	var wg sync.WaitGroup
	db, err := gorm.Open(sqlite.Open(cfg.Database.Url))
	if err != nil {
		return fmt.Errorf("[Start] gorm.Open: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if s, err := db.DB(); err != nil {
			return
		} else {
			_ = s.Close()
		}
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
