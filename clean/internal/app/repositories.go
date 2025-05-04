package app

import (
	"fmt"
	"log/slog"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

type repositories struct {
	chats       domain.ChatsRepository
	invitations domain.InvitationsRepository
	members     domain.MembersRepository
	users       domain.UsersRepository
	sessions    domain.SessionsRepository
	loginCreds  domain.LoginCredentialsRepository
}

func initSqliteRepositories(config sqlite.Config) (*repositories, func(), error) {
	factory, err := sqlite.InitRepositoryFactory(config)
	if err != nil {
		return nil, func() {}, fmt.Errorf("sqlite.InitRepositoryFactory: %w", err)
	}

	rs := &repositories{
		chats:       factory.NewChatsRepository(),
		invitations: factory.NewInvitationsRepository(),
		members:     factory.NewMembersRepository(),
		users:       factory.NewUsersRepository(),
		sessions:    factory.NewSessionsRepository(),
		loginCreds:  factory.NewLoginCredentialsRepository(),
	}

	return rs, func() {
		if err := factory.Close(); err != nil {
			slog.Error("initSqliteRepositories: factory.Close: " + err.Error())
		}
	}, nil
}
