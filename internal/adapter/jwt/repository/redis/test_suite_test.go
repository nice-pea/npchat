package redisCache_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	redisCache "github.com/nice-pea/npchat/internal/adapter/jwt/repository/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

type testSuite struct {
	suite.Suite
	cleanUp  func()
	DSN      string
	RedisCli redisCache.JWTIssuanceRegistry
}

func Test_TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

func (suite *testSuite) newRedisContainer() {
	ctx := context.Background()
	container, err := redisContainer.Run(ctx, "redis:8.2.1")
	suite.Require().NoError(err)
	dsn, err := container.ConnectionString(ctx)
	suite.Require().NoError(err)

	suite.cleanUp = func() {
		suite.Require().NotNil(container)
		container.Terminate(ctx)
	}

	suite.DSN = dsn

	redisCli, err := redisCache.Init(redisCache.Config{
		DSN: dsn,
	})

	suite.Require().NoError(err)
	suite.RedisCli = redisCache.JWTIssuanceRegistry{redisCli, 2 * time.Minute}

}

// SetupSuite выполняется один раз перед всеми тестами
func (suite *testSuite) SetupSuite() {
	suite.newRedisContainer()
}

// SetupTest выполняется перед каждым тестом
func (suite *testSuite) SetupTest() {

}

// TearDownTest выполняется после каждого теста
func (suite *testSuite) TearDownTest() {
	// println("TearDownTest executed")
}

func (suite *testSuite) SetupSubTest() {
	// Очищаем Redis перед каждым подтестом
	suite.RedisCli.FlushDB(context.Background())
	suite.RedisCli.Ttl = 2 * time.Minute
}

// TearDownSuite выполняется после всех тестов
func (suite *testSuite) TearDownSuite() {
	suite.cleanUp()
}

func (suite *testSuite) getIssueTime(sessionID uuid.UUID) time.Time {

	ctx := context.Background()
	var issueTime time.Time

	err := suite.RedisCli.Get(ctx, sessionID.String()).Scan(&issueTime)
	if !errors.Is(err, redis.Nil) {
		suite.Require().NoError(err)
	}
	return issueTime
}

func (suite *testSuite) setIssueTime(sessionID uuid.UUID, issueTime time.Time, ttl time.Duration) {
	ctx := context.Background()
	err := suite.RedisCli.Set(ctx, sessionID.String(), issueTime, ttl).Err()
	suite.Require().NoError(err)
}

func (suite *testSuite) redisEmpty() {
	ctx := context.Background()
	keys, err := suite.RedisCli.Keys(ctx, "*").Result()
	suite.Require().NoError(err)
	suite.Require().Empty(keys, "Redis должен быть пустым")
}

func (suite *testSuite) redisKeys() []string {
	ctx := context.Background()
	keys, err := suite.RedisCli.Keys(ctx, "*").Result()
	suite.Require().NoError(err)
	return keys
}
