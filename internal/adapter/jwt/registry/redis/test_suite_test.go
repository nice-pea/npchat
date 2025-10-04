package redisRegistry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	redisRegistry "github.com/nice-pea/npchat/internal/adapter/jwt/registry/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/suite"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

type testSuite struct {
	suite.Suite
	Terminate func()
	CleanUp   func()
	DSN       string
	RedisCli  redisRegistry.Registry
}

// Test_TestSuite запускает тестовый сценарий
func Test_TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

// newRedisContainer инициализирует Redis-контейнер для тестов
func (suite *testSuite) newRedisContainer() {
	ctx := context.Background()
	container, err := redisContainer.Run(ctx, "redis:8.2.1")
	suite.Require().NoError(err)
	dsn, err := container.ConnectionString(ctx)
	suite.Require().NoError(err)

	suite.Terminate = func() {
		suite.Require().NotNil(container)
		_ = container.Terminate(ctx)
	}
	suite.CleanUp = func() {
		status := suite.RedisCli.Client.FlushDB(context.Background())
		suite.Require().NoError(status.Err())
		suite.RedisCli.Ttl = 2 * time.Minute
	}
	suite.DSN = dsn

	redisCli, err := redisRegistry.Init(redisRegistry.Config{
		DSN: dsn,
	})

	suite.Require().NoError(err)
	suite.RedisCli = redisRegistry.Registry{redisCli, 2 * time.Minute}

}

// SetupSuite выполняется один раз перед всеми тестами
func (suite *testSuite) SetupSuite() {
	suite.newRedisContainer()
}

func (suite *testSuite) SetupSubTest() {
	// Очищаем Redis перед каждым подтестом
	suite.CleanUp()
}

// TearDownSuite выполняется после всех тестов
func (suite *testSuite) TearDownSuite() {
	suite.Terminate()
}

// getIssueTime возвращает время анулирования токена из Redis по sessionID
func (suite *testSuite) getIssueTime(sessionID uuid.UUID) time.Time {
	var issueTime time.Time
	err := suite.RedisCli.Client.Get(context.Background(), sessionID.String()).Scan(&issueTime)
	if !errors.Is(err, redis.Nil) {
		suite.Require().NoError(err)
	}
	return issueTime
}

// setIssueTime записывает время анулирования токена в Redis с указанным TTL
func (suite *testSuite) setIssueTime(sessionID uuid.UUID, issueTime time.Time, ttl time.Duration) {
	err := suite.RedisCli.Client.Set(context.Background(), sessionID.String(), issueTime, ttl).Err()
	suite.Require().NoError(err)
}

// requireIsRedisEmpty проверяет, что Redis пуст
func (suite *testSuite) requireIsRedisEmpty() {
	ctx := context.Background()
	keys, err := suite.RedisCli.Client.Keys(ctx, "*").Result()
	suite.Require().NoError(err)
	suite.Require().Empty(keys, "Redis должен быть пустым")
}

// redisKeys возвращает список всех ключей в Redis
func (suite *testSuite) redisKeys() []string {
	ctx := context.Background()
	keys, err := suite.RedisCli.Client.Keys(ctx, "*").Result()
	suite.Require().NoError(err)
	return keys
}
