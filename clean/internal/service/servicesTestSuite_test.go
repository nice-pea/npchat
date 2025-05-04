package service

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite"
)

type servicesTestSuite struct {
	suite.Suite
	factory *sqlite.RepositoryFactory
	rr      struct {
		chats            domain.ChatsRepository
		members          domain.MembersRepository
		invitations      domain.InvitationsRepository
		sessions         domain.SessionsRepository
		users            domain.UsersRepository
		loginCredentials domain.LoginCredentialsRepository
	}
	ss struct {
		chats            *Chats
		members          *Members
		invitations      *Invitations
		sessions         *Sessions
		loginCredentials *LoginCredentials
	}
}

func Test_ServicesTestSuite(t *testing.T) {
	suite.Run(t, new(servicesTestSuite))
}

// SetupSubTest выполняется перед каждым подтестом, связанным с suite
func (suite *servicesTestSuite) SetupSubTest() {
	var err error
	require := suite.Require()

	// Инициализация SQLiteMemory
	suite.factory, err = sqlite.InitRepositoryFactory(sqlite.Config{
		MigrationsDir: "../../migrations/repository/sqlite",
	})
	require.NoError(err)

	// Инициализация репозиториев
	suite.rr.chats = suite.factory.NewChatsRepository()
	suite.rr.members = suite.factory.NewMembersRepository()
	suite.rr.invitations = suite.factory.NewInvitationsRepository()
	suite.rr.users = suite.factory.NewUsersRepository()
	suite.rr.sessions = suite.factory.NewSessionsRepository()
	suite.rr.loginCredentials = suite.factory.NewLoginCredentialsRepository()

	// Создание сервисов
	suite.ss.chats = &Chats{
		ChatsRepo:   suite.rr.chats,
		MembersRepo: suite.rr.members,
	}
	suite.ss.members = &Members{
		ChatsRepo:   suite.rr.chats,
		MembersRepo: suite.rr.members,
	}
	suite.ss.invitations = &Invitations{
		ChatsRepo:       suite.rr.chats,
		MembersRepo:     suite.rr.members,
		InvitationsRepo: suite.rr.invitations,
		UsersRepo:       suite.rr.users,
	}
	suite.ss.sessions = &Sessions{
		SessionsRepo: suite.rr.sessions,
	}
	suite.ss.loginCredentials = &LoginCredentials{
		LoginCredentialsRepo: suite.rr.loginCredentials,
		SessionsRepo:         suite.rr.sessions,
	}
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *servicesTestSuite) TearDownSubTest() {
	err := suite.factory.Close()
	suite.Require().NoError(err)
}

// saveChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveChat(chat domain.Chat) domain.Chat {
	err := suite.rr.chats.Save(chat)
	suite.Require().NoError(err)

	return chat
}

// saveMember сохраняет участника в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveMember(member domain.Member) domain.Member {
	err := suite.rr.members.Save(member)
	suite.Require().NoError(err)

	return member
}

// saveInvitation сохраняет приглашение в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveInvitation(invitation domain.Invitation) domain.Invitation {
	err := suite.rr.invitations.Save(invitation)
	suite.Require().NoError(err)

	return invitation
}

// saveUser сохраняет пользователя в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveUser(user domain.User) domain.User {
	err := suite.rr.users.Save(user)
	suite.Require().NoError(err)

	return user
}
