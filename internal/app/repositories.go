package app

import (
	"fmt"
	"log/slog"

	"github.com/saime-0/nice-pea-chat/internal/domain/chatt"
	"github.com/saime-0/nice-pea-chat/internal/domain/sessionn"
	"github.com/saime-0/nice-pea-chat/internal/domain/userr"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

type repositories struct {
	chats    chatt.Repository
	users    userr.Repository
	sessions sessionn.Repository
}

func initSqliteRepositories(config sqlite.Config) (*repositories, func(), error) {
	factory, err := sqlite.InitRepositoryFactory(config)
	if err != nil {
		return nil, func() {}, fmt.Errorf("sqlite.InitRepositoryFactory: %w", err)
	}

	rs := &repositories{
		chats:    factory.NewChattRepository(),
		users:    factory.NewUserrRepository(),
		sessions: factory.NewSessionnRepository(),
	}

	return rs, func() {
		if err := factory.Close(); err != nil {
			slog.Error("initSqliteRepositories: factory.Close: " + err.Error())
		}
	}, nil
}
