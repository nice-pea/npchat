package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/saime-0/nice-pea-chat/internal/app/config"
	"github.com/saime-0/nice-pea-chat/internal/http"
	serviceL10nDB "github.com/saime-0/nice-pea-chat/internal/service/l10n/db"
)

func Start(ctx context.Context, cfg config.Config) (err error) {
	var wg sync.WaitGroup
	var mainDB, l10nDB *gorm.DB

	if mainDB, err = sqliteConnect(ctx, cfg.Database.DSN, &wg); err != nil {
		return fmt.Errorf("[Start] mainDB.sqliteConnect: %w", err)
	}
	if l10nDB, err = sqliteConnect(ctx, cfg.L10n.DSN, &wg); err != nil {
		return fmt.Errorf("[Start] l10n.sqliteConnect: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := (http.ServerParams{
			Ctx:  ctx,
			Addr: cfg.App.Address,
			L10n: &serviceL10nDB.Service{DB: l10nDB},
			DB:   mainDB,
		}.StartServer()); err != nil {
			log.Printf("[Start] http.StartServer: %s", err.Error())
		}
	}()

	log.Println("[Start] Receive ctx.Done, wait when components stop the work")
	wg.Wait()
	log.Println("[Start] Components done the work")

	return nil
}

func sqliteConnect(ctx context.Context, dsn string, wg *sync.WaitGroup) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("[Start] gorm.Open: %w", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		if s, err := db.DB(); err == nil {
			_ = s.Close()
		}
	}()
	return db, nil
}
