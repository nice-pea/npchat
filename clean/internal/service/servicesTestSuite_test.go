package service

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/saime-0/nice-pea-chat/internal/domain"
	"github.com/saime-0/nice-pea-chat/internal/repository/sqlite/memory"
)

type servicesTestSuite struct {
	suite.Suite
	sqliteMemory       *memory.SQLiteMemory
	chatsService       *Chats
	membersService     *Members
	invitationsService *Invitations
}

func Test_ServicesTestSuite(t *testing.T) {
	suite.Run(t, new(servicesTestSuite))
}

// SetupSubTest выполняется перед каждым подтестом, связанным с suite
func (suite *servicesTestSuite) SetupSubTest() {
	var err error
	require := suite.Require()

	// Инициализация SQLiteMemory
	suite.sqliteMemory, err = memory.Init(memory.Config{MigrationsDir: "../../migrations/repository/sqlite/memory"})
	require.NoError(err)

	// Инициализация репозиториев
	chatsRepository, err := suite.sqliteMemory.NewChatsRepository()
	require.NoError(err)
	membersRepository, err := suite.sqliteMemory.NewMembersRepository()
	require.NoError(err)
	invitationsRepository, err := suite.sqliteMemory.NewInvitationsRepository()
	require.NoError(err)
	usersRepository, err := suite.sqliteMemory.NewUsersRepository()
	require.NoError(err)

	// Создание сервисов
	suite.chatsService = &Chats{
		ChatsRepo:   chatsRepository,
		MembersRepo: membersRepository,
	}
	suite.membersService = &Members{
		ChatsRepo:   chatsRepository,
		MembersRepo: membersRepository,
	}
	suite.invitationsService = &Invitations{
		ChatsRepo:       chatsRepository,
		MembersRepo:     membersRepository,
		InvitationsRepo: invitationsRepository,
		UsersRepo:       usersRepository,
	}
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *servicesTestSuite) TearDownSubTest() {
	err := suite.sqliteMemory.Close()
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
