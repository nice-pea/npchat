package redisCache_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	redisCache "github.com/nice-pea/npchat/internal/adapter/jwt/repository/redis"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type testSuite struct {
	suite.Suite
	cleanUp     func()
	ExposedAddr string
	RedisCli    redisCache.JWTIssuanceRegistry
}

func Test_TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) newRedisContainer() {

	req := testcontainers.ContainerRequest{
		Image:        "redis:8.2.1",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("Ready to accept connections"),
	}

	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	suite.Require().NoError(err)

	suite.cleanUp = func() {
		suite.Require().NotNil(container)
		container.Terminate(ctx)
	}

	// Получаем адрес контейнера
	endpoint, err := container.Endpoint(ctx, "")
	suite.Require().NoError(err)

	suite.ExposedAddr = endpoint

	redisCli, err := redisCache.Init(redisCache.Config{
		Addr: suite.ExposedAddr,
	})

	suite.Require().NoError(err)
	suite.RedisCli = redisCache.JWTIssuanceRegistry{redisCli}

}

// SetupSuite выполняется один раз перед всеми тестами
func (suite *testSuite) SetupSuite() {
	// suite.newRedisContainer()
}

// SetupTest выполняется перед каждым тестом
func (suite *testSuite) SetupTest() {
	suite.newRedisContainer()

}

// TearDownTest выполняется после каждого теста
func (suite *testSuite) TearDownTest() {
	println("TearDownTest executed")
}

// TearDownSuite выполняется после всех тестов
func (suite *testSuite) TearDownSuite() {
	suite.cleanUp()
}

func (suite *testSuite) getIssueTime(sessionID uuid.UUID) time.Time {

	ctx := context.Background()
	cmd := suite.RedisCli.Get(ctx, sessionID.String())
	res, err := cmd.Result()
	suite.Assert().NoError(err)

	issueTime, err := time.Parse(time.Layout, res)
	suite.Assert().NoError(err)

	return issueTime
}
