package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os/exec"

	controller "github.com/saime-0/nice-pea-chat/internal/controller/http"
	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

type repositories struct {
	chats       domain.ChatsRepository
	invitations domain.InvitationsRepository
	members     domain.MembersRepository
	users       domain.UsersRepository
	sessions    domain.SessionsRepository
}

type services struct {
	chats       service.Chats
	invitations service.Invitations
	members     service.Members
	sessions    service.Sessions
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

func initServices(repos *repositories) *services {
	return &services{
		chats: service.Chats{
			ChatsRepo:   repos.chats,
			MembersRepo: repos.members,
		},
		invitations: service.Invitations{
			ChatsRepo:       repos.chats,
			MembersRepo:     repos.members,
			InvitationsRepo: repos.invitations,
			UsersRepo:       repos.users,
		},
		members: service.Members{
			MembersRepo: repos.members,
			ChatsRepo:   repos.chats,
		},
		sessions: service.Sessions{
			SessionsRepo: repos.sessions,
		},
	}
}

func Run(ctx context.Context) error {
	// Execute pwd command and print output
	out, err := exec.Command("pwd").Output()
	if err != nil {
		return fmt.Errorf("failed to execute pwd command: %w", err)
	}
	fmt.Printf("Current working directory: %s", out)

	// Инициализация репозиториев
	repos, closer, err := initSqliteRepositories(sqlite.Config{
		MigrationsDir: "migrations/repository/sqlite",
	})
	if err != nil {
		return err
	}
	defer closer()

	// Инициализация сервисов
	ss := initServices(repos)

	// Инициализация контроллера
	ctrl := controller.InitController(ss.chats, ss.invitations, ss.members)

	// Запуск сервера
	server := &http.Server{
		Addr:    ":8080",
		Handler: ctrl,
	}

	errChan := make(chan error, 1)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err = <-errChan:
		return err
	case <-ctx.Done():
		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("server shutdown: %w", err)
		}
		return nil
	}
}
