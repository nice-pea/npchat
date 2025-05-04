package http

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
	"github.com/saime-0/nice-pea-chat/internal/service"
)

type controllerTestSuite struct {
	suite.Suite
	ctrl    *Controller
	factory *sqlite.RepositoryFactory
	rr      struct {
		chats       domain.ChatsRepository
		members     domain.MembersRepository
		invitations domain.InvitationsRepository
		sessions    domain.SessionsRepository
		users       domain.UsersRepository
		loginCreds  domain.LoginCredentialsRepository
	}
	ss struct {
		chats            *service.Chats
		members          *service.Members
		invitations      *service.Invitations
		sessions         *service.Sessions
		loginCredentials *service.LoginCredentials
	}
	server *httptest.Server
}

func Test_ServicesTestSuite(t *testing.T) {
	suite.Run(t, new(controllerTestSuite))
}

// SetupSubTest выполняется перед каждым подтестом, связанным с suite
func (suite *controllerTestSuite) SetupSubTest() {
	var err error
	require := suite.Require()

	// Инициализация SQLiteMemory
	suite.factory, err = sqlite.InitRepositoryFactory(sqlite.Config{
		MigrationsDir: "../../../migrations/repository/sqlite",
	})
	require.NoError(err)

	// Инициализация репозиториев
	suite.rr.chats = suite.factory.NewChatsRepository()
	suite.rr.members = suite.factory.NewMembersRepository()
	suite.rr.invitations = suite.factory.NewInvitationsRepository()
	suite.rr.users = suite.factory.NewUsersRepository()
	suite.rr.sessions = suite.factory.NewSessionsRepository()
	suite.rr.loginCreds = suite.factory.NewLoginCredentialsRepository()

	// Создание сервисов
	suite.ss.chats = &service.Chats{
		ChatsRepo:   suite.rr.chats,
		MembersRepo: suite.rr.members,
	}
	suite.ss.members = &service.Members{
		ChatsRepo:   suite.rr.chats,
		MembersRepo: suite.rr.members,
	}
	suite.ss.invitations = &service.Invitations{
		ChatsRepo:       suite.rr.chats,
		MembersRepo:     suite.rr.members,
		InvitationsRepo: suite.rr.invitations,
		UsersRepo:       suite.rr.users,
	}
	suite.ss.sessions = &service.Sessions{
		SessionsRepo: suite.rr.sessions,
	}
	suite.ss.loginCredentials = &service.LoginCredentials{
		LoginCredentialsRepo: suite.rr.loginCreds,
		SessionsRepo:         suite.rr.sessions,
	}

	suite.ctrl = &Controller{
		chats:            suite.ss.chats,
		invitations:      suite.ss.invitations,
		members:          suite.ss.members,
		sessions:         suite.ss.sessions,
		loginCredentials: suite.ss.loginCredentials,
		ServeMux:         http.ServeMux{},
	}
	suite.ctrl.registerHandlers()
	suite.server = httptest.NewServer(suite.ctrl)
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *controllerTestSuite) TearDownSubTest() {
	suite.ctrl = nil
	suite.server.Close() //nolint:errcheck
	err := suite.factory.Close()
	suite.Require().NoError(err)
}
