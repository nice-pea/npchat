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
		chats       domain.ChatsRepository
		members     domain.MembersRepository
		invitations domain.InvitationsRepository
		sessions    domain.SessionsRepository
		users       domain.UsersRepository
	}
	ss struct {
		chats       *Chats
		members     *Members
		invitations *Invitations
		sessions    *Sessions
	}
	chatsService       *Chats
	membersService     *Members
	invitationsService *Invitations
	sessionsService    *Sessions
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
	suite.rr.chats, err = suite.factory.NewChatsRepository()
	require.NoError(err)
	suite.rr.members, err = suite.factory.NewMembersRepository()
	require.NoError(err)
	suite.rr.invitations, err = suite.factory.NewInvitationsRepository()
	require.NoError(err)
	suite.rr.users, err = suite.factory.NewUsersRepository()
	require.NoError(err)
	suite.rr.sessions, err = suite.factory.NewSessionsRepository()
	require.NoError(err)

	// Создание сервисов
	suite.chatsService = &Chats{
		ChatsRepo:   suite.rr.chats,
		MembersRepo: suite.rr.members,
	}
	suite.membersService = &Members{
		ChatsRepo:   suite.rr.chats,
		MembersRepo: suite.rr.members,
	}
	suite.invitationsService = &Invitations{
		ChatsRepo:       suite.rr.chats,
		MembersRepo:     suite.rr.members,
		InvitationsRepo: suite.rr.invitations,
		UsersRepo:       suite.rr.users,
	}

	suite.ss.chats = suite.chatsService
	suite.ss.members = suite.membersService
	suite.ss.invitations = suite.invitationsService
	suite.ss.sessions = &Sessions{
		SessionsRepo: suite.rr.sessions,
	}
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *servicesTestSuite) TearDownSubTest() {
	err := suite.factory.Close()
	suite.Require().NoError(err)
}

// saveChat сохраняет чат в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveChat(chat domain.Chat) domain.Chat {
	err := suite.chatsService.ChatsRepo.Save(chat)
	suite.Require().NoError(err)

	return chat
}

// saveMember сохраняет участника в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveMember(member domain.Member) domain.Member {
	err := suite.membersService.MembersRepo.Save(member)
	suite.Require().NoError(err)

	return member
}

// saveInvitation сохраняет приглашение в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveInvitation(invitation domain.Invitation) domain.Invitation {
	err := suite.invitationsService.InvitationsRepo.Save(invitation)
	suite.Require().NoError(err)

	return invitation
}

// saveUser сохраняет пользователя в репозиторий, в случае ошибки завершит тест
func (suite *servicesTestSuite) saveUser(user domain.User) domain.User {
	err := suite.invitationsService.UsersRepo.Save(user)
	suite.Require().NoError(err)

	return user
}
