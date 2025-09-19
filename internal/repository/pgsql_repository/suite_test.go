package pgsqlRepository

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	testifySuite "github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/nice-pea/npchat/internal/domain/chatt"
	"github.com/nice-pea/npchat/internal/domain/sessionn"
	"github.com/nice-pea/npchat/internal/domain/userr"
)

type Suite struct {
	testifySuite.Suite
	factory       *Factory
	factoryCloser func()
	RR            struct {
		//Chats    *ChattRepository
		//Sessions *SessionnRepository
		//Users    *UserrRepository
		Chats    chatt.Repository
		Sessions sessionn.Repository
		Users    userr.Repository
	}
}

func Test_TestSuite(t *testing.T) {
	testifySuite.Run(t, new(Suite))
}

var (
	pgsqlDSN = os.Getenv("TEST_PGSQL_DSN")
)

// newPgsqlExternalFactory создает фабрику репозиториев для тестирования, реализованных с помощью подключения к postgres по DSN
func (suite *Suite) newPgsqlExternalFactory(dsn string) (*Factory, func()) {
	factory, err := InitFactory(Config{
		DSN: dsn,
	})
	suite.Require().NoError(err)

	return factory, func() { _ = factory.Close() }
}

// newPgsqlContainerFactory создает фабрику репозиториев для тестирования, реализованных с помощью postgres контейнеров
func (suite *Suite) newPgsqlContainerFactory() (f *Factory, closer func()) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Поиск скриптов с миграциями
	_, b, _, _ := runtime.Caller(0)
	migrationsDir := filepath.Join(filepath.Dir(b), "../../../infra/pgsql/init/*.up.sql")
	migrations, err := filepath.Glob(migrationsDir)
	suite.Require().NoError(err)
	suite.Require().NotZero(migrations)

	postgresContainer, err := postgres.Run(ctx,
		"postgres:17",
		postgres.WithInitScripts(migrations...),
		postgres.WithDatabase("test_npc_db"),
		postgres.WithUsername("test_npc_user"),
		postgres.WithPassword("test_npc_password"),
		postgres.BasicWaitStrategies(),
	)
	suite.Require().NoError(err)
	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	suite.Require().NoError(err)

	f, err = InitFactory(Config{
		DSN: dsn,
	})
	suite.Require().NoError(err)

	return f, func() {
		_ = f.Close()
		_ = postgresContainer.Terminate(context.Background())
	}
}

// SetupTest выполняется перед каждым тестом, связанным с suite
func (suite *Suite) SetupTest() {
	// Инициализация фабрики репозиториев
	if pgsqlDSN != "" {
		suite.factory, suite.factoryCloser = suite.newPgsqlExternalFactory(pgsqlDSN)
	} else {
		suite.factory, suite.factoryCloser = suite.newPgsqlContainerFactory()
	}

	// Инициализация репозиториев
	suite.RR.Chats = suite.factory.NewChattRepository()
	suite.RR.Users = suite.factory.NewUserrRepository()
	suite.RR.Sessions = suite.factory.NewSessionnRepository()
}

// TearDownSubTest выполняется после каждого подтеста, связанного с suite
func (suite *Suite) TearDownSubTest() {
	err := suite.factory.Cleanup()
	suite.Require().NoError(err)
}

// TearDownTest выполняется после каждого подтеста, связанного с suite
func (suite *Suite) TearDownTest() {
	suite.factoryCloser()
}
