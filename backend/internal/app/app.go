package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/saime-0/nice-pea-chat/internal/config"
	"github.com/saime-0/nice-pea-chat/internal/http"
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
			L10n: l10n{},
			DB:   db,
		}.StartServer()); err != nil {
			log.Printf("[Start] http.StartServer: %s", err.Error())
		}
	}()

	log.Println("[Start] Receive ctx.Done, wait when components stop the work")
	wg.Wait()
	log.Println("[Start] Components done the work")
	return nil
}

type l10n struct{}

func (l l10n) Localize(code, locale string, vars map[string]string) (string, error) {
	switch code {
	case "none:ok":
		return "ok", nil
	default:
		return "", fmt.Errorf("unknown code: %s", code)
	}
}
