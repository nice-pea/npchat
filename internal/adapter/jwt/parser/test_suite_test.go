package jwtParser

import (
	"context"
	"testing"

	jwt2 "github.com/cristalhq/jwt/v5"
	"github.com/google/uuid"
	"github.com/nice-pea/npchat/internal/adapter/jwt"
	redisCache "github.com/nice-pea/npchat/internal/adapter/jwt/repository/redis"
	"github.com/stretchr/testify/suite"
	redisContainer "github.com/testcontainers/testcontainers-go/modules/redis"
)

type testSuite struct {
	suite.Suite
	cleanUp func()
	cfg     jwt.Config
	Parser  Parser
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

	suite.cfg = jwt.Config{
		SecretKey:                     "secret",
		VerifyTokenWithAdvancedChecks: true,
		RedisDSN:                      dsn,
	}
	cli, err := redisCache.Init(redisCache.Config{
		DSN: dsn,
	})
	suite.Require().NoError(err)

	suite.Parser = Parser{
		Config:                        suite.cfg,
		VerifyTokenWithAdvancedChecks: true,
		cache:                         redisCache.JWTIssuanceRegistry{Client: cli},
	}
}

func Test_TestSuite(t *testing.T) {
	suite.Run(t, new(testSuite))
}

// SetupSuite выполняется один раз перед всеми тестами
func (suite *testSuite) SetupSuite() {
	suite.newRedisContainer()
}

func (suite *testSuite) SetupSubTest() {
	// Очищаем Redis перед каждым подтестом
	if suite.Parser.cache.Client != nil {
		suite.Parser.cache.FlushDB(context.Background())
	}
	cli, err := redisCache.Init(redisCache.Config{
		DSN: suite.cfg.RedisDSN,
	})
	suite.Require().NoError(err)
	suite.Parser = Parser{
		Config:                        suite.cfg,
		VerifyTokenWithAdvancedChecks: true,
		cache:                         redisCache.JWTIssuanceRegistry{Client: cli},
	}

}

// TearDownSuite выполняется после всех тестов
func (suite *testSuite) TearDownSuite() {
	suite.cleanUp()
}

func (suite *testSuite) mustParseUuid(s string) uuid.UUID {
	uid, err := uuid.Parse(s)
	suite.Require().NoError(err)
	return uid
}

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
