package jwtParser

import (
	"context"
	"testing"

	jwt2 "github.com/cristalhq/jwt/v5"
	"github.com/nice-pea/npchat/internal/adapter/jwt"
	redisRegistry "github.com/nice-pea/npchat/internal/adapter/jwt/registry/redis"
	"github.com/stretchr/testify/suite"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

type testSuite struct {
	suite.Suite
	CleanUp   func()
	Terminate func()
	cfg       jwt.Config
	Parser    Parser
}

// newRedisContainer настраивает Redis контейнер для тестов
func (suite *testSuite) newRedisContainer() *redisContainer.RedisContainer {
	ctx := context.Background()
	container, err := redisContainer.Run(ctx, "redis:8.2.1")
	suite.Require().NoError(err)

	suite.Terminate = func() {
		suite.Require().NotNil(container)
		_ = container.Terminate(ctx)
	}

	return container
}

// Test_TestSuite запускает набор тестов
func Test_TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

// SetupSuite выполняется один раз перед всеми тестами
func (suite *testSuite) SetupSuite() {
	// создание контейнера
	container := suite.newRedisContainer()

	// создание конфигуарции с Registry
	dsn, err := container.ConnectionString(context.Background())
	suite.Require().NoError(err)
	suite.cfg = jwt.Config{
		SecretKey:                   "secret",
		VerifyTokenWithInvalidation: true,
		RedisDSN:                    dsn,
	}
}

// SetupSubTest выполняется перед каждым подтестом
func (suite *testSuite) SetupSubTest() {
	// Очищаем Redis перед каждым подтестом
	if suite.CleanUp != nil {
		suite.CleanUp()
	}

	// Пересоздаем Parser
	cli, err := redisRegistry.Init(redisRegistry.Config{
		DSN: suite.cfg.RedisDSN,
	})
	suite.Require().NoError(err)
	suite.cfg.VerifyTokenWithInvalidation = true
	registry := &redisRegistry.Registry{Client: cli}
	suite.Parser = Parser{
		Config:   suite.cfg,
		Registry: registry,
	}
	suite.CleanUp = func() {
		if suite.Parser.Registry != nil {
			status := registry.Client.FlushDB(context.Background())
			suite.Require().NoError(status.Err())
		}
	}
}

// TearDownSuite выполняется после всех тестов
func (suite *testSuite) TearDownSuite() {
	suite.Terminate()
}

// createJWT создает JWT токен для тестов
func (suite *testSuite) createJWT(secret string, claims map[string]any) string {
	// создаем Signer
	signer, err := jwt2.NewSignerHS(jwt2.HS256, []byte(secret))
	suite.Require().NoError(err)
	// создаем Builder
	builder := jwt2.NewBuilder(signer)

	// создаем токен
	token, err := builder.Build(claims)
	suite.Require().NoError(err)

	return token.String()
}

// parserWithOutRegistry создает парсер без Registry
func (suite *testSuite) parserWithOutRegistry(secret string) Parser {
	return Parser{Config: jwt.Config{SecretKey: secret}}
}
