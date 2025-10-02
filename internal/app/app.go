package app

import (
	"context"
	"fmt"
	"log/slog"

	"golang.org/x/sync/errgroup"

	"github.com/nice-pea/npchat/internal/common"
	"github.com/nice-pea/npchat/internal/controller/http2"
)

func Run(ctx context.Context, cfg Config, buildInfo common.BuildInfo) error {
	g, ctx := errgroup.WithContext(ctx)

	slog.SetLogLoggerLevel(slogLevel(cfg.LogLevel))
	slog.Info(fmt.Sprintf("Уровень логирования: %s", cfg.LogLevel))

	// Инициализация репозиториев
	rr, closeRepos, err := initPgsqlRepositories(cfg.Pgsql)
	if err != nil {
		return err
	}
	defer closeRepos()

	// Инициализация адаптеров
	aa := initAdapters(cfg)

	// Инициализация сервисов
	uc := initUsecases(rr, aa)

	// Инициализация и Запуск http контроллера
	g.Go(func() error {
		return http2.RunHttpServer(ctx, uc, aa.eventBus, aa.jwtIssuer, aa.jwtParser, cfg.Http2, buildInfo)
	})

	return g.Wait()
}
