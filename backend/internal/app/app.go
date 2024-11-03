package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/jackc/pgx/v5"
	"gorm.io/driver/postgres"
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
	var pgxConn *pgx.Conn
	if mainDB, err = postgresConnect(ctx, cfg.Database.DSN, &wg); err != nil {
		return fmt.Errorf("[Start] mainDB.postgresConnect: %w", err)
	}
	if pgxConn, err = pgxConnect(ctx, cfg.Database.DSN, &wg); err != nil {
		return fmt.Errorf("[Start] mainDB.postgresConnect: %w", err)
	}
	if l10nDB, err = sqliteConnect(ctx, cfg.L10n.DSN, &wg); err != nil {
		return fmt.Errorf("[Start] l10n.sqliteConnect: %w", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := (http.ServerParams{
			Ctx:     ctx,
			Addr:    cfg.App.Address,
			L10n:    &serviceL10nDB.Service{DB: l10nDB},
			DB:      mainDB,
			PGXconn: pgxConn,
		}.StartServer()); err != nil {
			log.Printf("[Start] http.StartServer: %s", err.Error())
		}
	}()

	log.Println("[Start] Receive ctx.Done, wait when components stop the work")
	wg.Wait()
	log.Println("[Start] Components done the work")

	return nil
}

func pgxConnect(ctx context.Context, dsn string, wg *sync.WaitGroup) (*pgx.Conn, error) {
	parseConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("[pgxConnect] ParseConfig: %w", err)
	}
	parseConfig.Tracer = &MyTracer{}
	connect, err := pgx.ConnectConfig(context.Background(), parseConfig)
	if err != nil {
		return nil, fmt.Errorf("[Start] pgx.Connect: %w", err)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()
		connect.Close(context.Background())
	}()
	return connect, nil
}

type MyTracer struct{}

func (t *MyTracer) TraceQueryStart(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryStartData) context.Context {
	log.Printf("Trace <\\\nsql: %s\nargs: %s\nTrace />", data.SQL, data.Args)
	return ctx
}

func (t *MyTracer) TraceQueryEnd(ctx context.Context, conn *pgx.Conn, data pgx.TraceQueryEndData) {
	log.Printf("EndTrace <\\\ntag: %s\nerr: %s\nEndTrace />", data.CommandTag, data.Err)
}

func postgresConnect(ctx context.Context, dsn string, wg *sync.WaitGroup) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
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
