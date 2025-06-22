package app

import (
	"fmt"
	"log/slog"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
	pgsqlRepository "github.com/nice-pea/npchat/internal/repository/pgsql_repository"
)

type repositories struct {
	chats    chatt.Repository
	users    userr.Repository
	sessions sessionn.Repository
}

//func initSqliteRepositories(config sqlite.Config) (*repositories, func(), error) {
//	factory, err := sqlite.InitRepositoryFactory(config)
//	if err != nil {
//		return nil, func() {}, fmt.Errorf("sqlite.InitRepositoryFactory: %w", err)
//	}
//
//	rs := &repositories{
//		chats:    factory.NewChattRepository(),
//		users:    factory.NewUserrRepository(),
//		sessions: factory.NewSessionnRepository(),
//	}
//
//	return rs, func() {
//		if err := factory.Close(); err != nil {
//			slog.Error("initSqliteRepositories: factory.Close: " + err.Error())
//		}
//	}, nil
//}

func initPgsqlRepositories(cfg pgsqlRepository.Config) (*repositories, func(), error) {
	factory, err := pgsqlRepository.InitFactory(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("pgsqlRepository.InitFactory: %w", err)
	}

	rs := &repositories{
		chats:    factory.NewChattRepository(),
		users:    factory.NewUserrRepository(),
		sessions: factory.NewSessionnRepository(),
	}

	closer := func() {
		if err := factory.Close(); err != nil {
			slog.Error("Закрыть соединение с pgsql: factory.Close: " + err.Error())
		}
	}

	return rs, closer, nil
}
