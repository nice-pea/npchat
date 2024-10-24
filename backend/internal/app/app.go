package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/app/config"
	"github.com/saime-0/nice-pea-chat/internal/http"
	l10nDB "github.com/saime-0/nice-pea-chat/internal/service/l10n/db"
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
		if s, err := db.DB(); err == nil {
			_ = s.Close()
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := (http.ServerParams{
			Ctx:  ctx,
			Addr: cfg.App.Address,
			L10n: &l10nDB.Service{
				DB: db,
			},
			DB: db,
		}.StartServer()); err != nil {
			log.Printf("[Start] http.StartServer: %s", err.Error())
		}
	}()

	log.Println("[Start] Receive ctx.Done, wait when components stop the work")
	wg.Wait()
	log.Println("[Start] Components done the work")
	return nil
}
