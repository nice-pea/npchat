package http

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type servicesTestSuite struct {
	suite.Suite
	ctrl Controller
	//sqliteMemory       *memory.SQLiteMemory
	//chatsService       *services.Chats
	//membersService     *services.Members
	//invitationsService *services.Invitations
}

func Test_ServicesTestSuite(t *testing.T) {
	suite.Run(t, new(servicesTestSuite))
}

// SetupSubTest выполняется перед каждым подтестом, связанным с suite
func (suite *servicesTestSuite) SetupSubTest() {
	//var err error
	//require := suite.Require()
	//
	//// Инициализация SQLiteMemory
	//suite.sqliteMemory, err = memory.Init(memory.Config{MigrationsDir: "../../migrations/repository/sqlite/memory"})
	//require.NoError(err)
	//
	//// Инициализация репозиториев
	//chatsRepository, err := suite.sqliteMemory.NewChatsRepository()
	//require.NoError(err)
	//membersRepository, err := suite.sqliteMemory.NewMembersRepository()
	//require.NoError(err)
	//invitationsRepository, err := suite.sqliteMemory.NewInvitationsRepository()
	//require.NoError(err)
	//usersRepository, err := suite.sqliteMemory.NewUsersRepository()
	//require.NoError(err)
	//
	//// Создание сервисов
	//suite.chatsService = &Chats{
	//	ChatsRepo:   chatsRepository,
	//	MembersRepo: membersRepository,
	//}
	//suite.membersService = &Members{
	//	ChatsRepo:   chatsRepository,
	//	MembersRepo: membersRepository,
	//}
	//suite.invitationsService = &Invitations{
	//	ChatsRepo:       chatsRepository,
	//	MembersRepo:     membersRepository,
	//	InvitationsRepo: invitationsRepository,
	//	UsersRepo:       usersRepository,
	//}
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *servicesTestSuite) TearDownSubTest() {
	//err := suite.sqliteMemory.Close()
	//suite.Require().NoError(err)
}
