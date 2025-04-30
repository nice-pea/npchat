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
}

func initSqliteRepositories(config sqlite.Config) (*repositories, func(), error) {
	factory, err := sqlite.InitRepositoryFactory(config)
	if err != nil {
		return nil, func() {}, fmt.Errorf("sqlite.InitRepositoryFactory: %w", err)
	}

	rs := new(repositories)

	if rs.chats, err = factory.NewChatsRepository(); err != nil {
		return nil, func() {}, fmt.Errorf("factory.NewChatsRepository: %w", err)
	}
	if rs.invitations, err = factory.NewInvitationsRepository(); err != nil {
		return nil, func() {}, fmt.Errorf("factory.NewInvitationsRepository: %w", err)
	}
	if rs.members, err = factory.NewMembersRepository(); err != nil {
		return nil, func() {}, fmt.Errorf("factory.NewMembersRepository: %w", err)
	}
	if rs.users, err = factory.NewUsersRepository(); err != nil {
		return nil, func() {}, fmt.Errorf("factory.NewUsersRepository: %w", err)
	}
	if rs.sessions, err = factory.NewSessionsRepository(); err != nil {
		return nil, func() {}, fmt.Errorf("factory.NewSessionsRepository: %w", err)
	}

	return rs, func() {
		if err := factory.Close(); err != nil {
			slog.Error("initSqliteRepositories: factory.Close: " + err.Error())
		}
	}, nil
}
